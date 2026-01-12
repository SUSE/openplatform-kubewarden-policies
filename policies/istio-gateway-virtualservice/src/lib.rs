use lazy_static::lazy_static;
use std::collections::HashSet;

use guest::prelude::*;
use kubewarden_policy_sdk::wapc_guest as guest;
extern crate kubewarden_policy_sdk as kubewarden;
use kubewarden::request::ValidationRequest;
use kubewarden::{logging, protocol_version_guest, validate_settings};

mod settings;
use settings::Settings;

use slog::{info, o, warn, Logger};
use thiserror::Error;

lazy_static! {
    static ref LOG_DRAIN: Logger = Logger::root(
        logging::KubewardenDrain::new(),
        o!("policy" => "istio-gateway")
    );
}
use kcr_networking_istio_io::v1::virtualservices::VirtualService;

#[no_mangle]
pub extern "C" fn wapc_init() {
    register_function("validate", validate);
    register_function("validate_settings", validate_settings::<Settings>);
    register_function("protocol_version", protocol_version_guest);
}

#[derive(Debug, Error)]
pub enum ValidateError {
    #[error("could not parse object: {error:?}")]
    ParseError { error: String },
    #[error("VirtualService spec must specify at least one gateway")]
    NoGatewayError,
    #[error("VirtualService spec must specify at least one host")]
    NoHostError,
    #[error("VirtualService must have a namespace to be evaluated")]
    NoNamespaceError,
    #[error(
        "Namespace '{namespace:?}' is not allowed to use the restricted gateway '{gateway:?}'"
    )]
    NamespaceNotAllowedError { namespace: String, gateway: String },
    #[error(
        "Host '{host:?}' is not allowed in namespace '{namespace:?}' for gateway '{gateway:?}'"
    )]
    HostNotAllowedError {
        host: String,
        namespace: String,
        gateway: String,
    },
    #[error("Service '{service:?}' (Port: {port}, Proto: {protocol}) is not allowed in namespace '{namespace:?}' for gateway '{gateway:?}'")]
    ServiceNotAllowedError {
        service: String,
        port: i32,
        protocol: String,
        namespace: String,
        gateway: String,
    },
}

static WILDCARD_HOSTS: &str = "*";
///
/// Helper struct to hold parsed destination data
#[derive(Debug, Clone, Hash, Eq, PartialEq)]
struct RequestedDestination {
    host: String,
    port: i32,
    protocol: String,
}

fn validate_request(req: ValidationRequest<Settings>) -> Result<(), ValidateError> {
    // Make sure to validate only istio's VirtualService objects
    if req.request.kind.kind != "VirtualService" {
        warn!(LOG_DRAIN, "Policy validates VirtualService only. Accepting resource"; "kind" => &req.request.kind.kind);
        return Ok(());
    }

    // Deserialize the incoming Kubernetes object into our VirtualService struct.
    let virtual_service =
        serde_json::from_value::<VirtualService>(req.request.object).map_err(|e| {
            ValidateError::ParseError {
                error: e.to_string(),
            }
        })?;

    let requested_gateways = virtual_service
        .spec
        .gateways
        .ok_or(ValidateError::NoGatewayError)?;

    // Get the namespace of the VirtualService. It must have one.
    let vs_namespace = virtual_service
        .metadata
        .namespace
        .ok_or(ValidateError::NoNamespaceError)?;

    let requested_hosts = virtual_service
        .spec
        .hosts
        .filter(|h| !h.is_empty())
        .unwrap_or_else(|| vec![WILDCARD_HOSTS.to_string()]);

    let mut all_requested_destinations: HashSet<RequestedDestination> = HashSet::new();

    // 1. HTTP Routes
    if let Some(routes) = &virtual_service.spec.http {
        for route in routes {
            if let Some(destinations) = &route.route {
                for dest in destinations {
                    let host = dest.destination.host.clone();
                    // Extract port number if present, otherwise 0
                    let port = dest
                        .destination
                        .port
                        .as_ref()
                        .map(|p| p.number.unwrap_or(0) as i32)
                        .unwrap_or(0);

                    all_requested_destinations.insert(RequestedDestination {
                        host,
                        port,
                        protocol: "HTTP".to_string(),
                    });
                }
            }
        }
    }

    // 2. TCP Routes
    if let Some(routes) = &virtual_service.spec.tcp {
        for route in routes {
            if let Some(destinations) = &route.route {
                for dest in destinations {
                    let host = dest.destination.host.clone();
                    let port = dest
                        .destination
                        .port
                        .as_ref()
                        .map(|p| p.number.unwrap_or(0) as i32)
                        .unwrap_or(0);

                    all_requested_destinations.insert(RequestedDestination {
                        host,
                        port,
                        protocol: "TCP".to_string(),
                    });
                }
            }
        }
    }

    // 3. TLS Routes
    if let Some(routes) = &virtual_service.spec.tls {
        for route in routes {
            if let Some(destinations) = &route.route {
                for dest in destinations {
                    let host = dest.destination.host.clone();
                    let port = dest
                        .destination
                        .port
                        .as_ref()
                        .map(|p| p.number.unwrap_or(0) as i32)
                        .unwrap_or(0);

                    all_requested_destinations.insert(RequestedDestination {
                        host,
                        port,
                        protocol: "TLS".to_string(),
                    });
                }
            }
        }
    }

    // The main validation logic.
    requested_gateways.iter().try_for_each(|gateway_name| {
        // Check if the gateway is in our list of restricted gateways.
        if let Some(namespaces_map) = req.settings.gateway_restrictions.get(gateway_name) {
            if let Some(rules_for_ns) = namespaces_map.get(&vs_namespace) {
                // --- Hostname Validation ---
                let allowed_hostnames: HashSet<_> = rules_for_ns
                    .iter()
                    .flat_map(|r| &r.hostnames)
                    .map(String::as_str)
                    .collect();

                if !allowed_hostnames.contains(WILDCARD_HOSTS) {
                    if let Some(unallowed_host) = requested_hosts
                        .iter()
                        .find(|&host| !allowed_hostnames.contains(host.as_str()))
                    {
                        return Err(ValidateError::HostNotAllowedError {
                            host: unallowed_host.clone(),
                            namespace: vs_namespace.clone(),
                            gateway: gateway_name.clone(),
                        });
                    }
                }

                // Service/Destination Validation
                // We must ensure that *every* requested destination matches *at least one* rule

                for req_dest in &all_requested_destinations {
                    let mut matches_rule = false;

                    for rule in rules_for_ns {
                        // Check 1: Host match
                        let host_matches = rule.destination_hosts.iter().any(|allowed_host| {
                            allowed_host == WILDCARD_HOSTS || allowed_host == &req_dest.host
                        });

                        if !host_matches {
                            continue;
                        }

                        // Check 2: Port match (0 in rule means ANY port)
                        let port_matches = rule.port == 0 || rule.port == req_dest.port;
                        if !port_matches {
                            continue;
                        }

                        // Check 3: Protocol match (Empty string in rule means ANY protocol)
                        let proto_matches = rule.protocol.is_empty()
                            || rule.protocol.eq_ignore_ascii_case(&req_dest.protocol);
                        if !proto_matches {
                            continue;
                        }

                        matches_rule = true;
                        break; // Found a matching rule for this destination, move to next requested dest
                    }

                    if !matches_rule {
                        return Err(ValidateError::ServiceNotAllowedError {
                            service: req_dest.host.clone(),
                            port: req_dest.port,
                            protocol: req_dest.protocol.clone(),
                            namespace: vs_namespace.clone(),
                            gateway: gateway_name.clone(),
                        });
                    }
                }

                Ok(()) // Both host and service validation passed for this gateway
            } else {
                Err(ValidateError::NamespaceNotAllowedError {
                    namespace: vs_namespace.clone(),
                    gateway: gateway_name.clone(),
                })
            }
        } else {
            // If the gateway is not in the settings, it's not restricted. Pass.
            Ok(())
        }
    })?;

    Ok(())
}

fn validate(payload: &[u8]) -> CallResult {
    let validation_req = ValidationRequest::<Settings>::new(payload)?;

    info!(LOG_DRAIN, "starting validation");
    match validate_request(validation_req) {
        Ok(_) => kubewarden::accept_request(),
        Err(err) => kubewarden::reject_request(Some(err.to_string()), None, None, None),
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use kubewarden_policy_sdk::test::Testcase;
    use settings::{Namespaces, Restriction, Settings};
    use std::collections::HashMap;

    static HTTP_PROTOCOL: &str = "http";
    static GW_01: &str = "gateway-01";
    static GW_02: &str = "gateway-02";
    static NS_01: &str = "ns-01";
    static NS_02: &str = "ns-02";
    static NS_04: &str = "ns-04";

    /// Helper function to build the policy settings
    fn build_settings() -> Settings {
        // Define restrictions for ns-01
        let ns_01_restrictions = vec![Restriction {
            hostnames: vec!["host-a.com".to_string(), "host-b.com".to_string()],
            destination_hosts: vec![WILDCARD_HOSTS.to_string()],
            port: 80,
            protocol: HTTP_PROTOCOL.to_string(),
        }];

        // Define restrictions for ns-02
        let ns_02_restrictions = vec![Restriction {
            hostnames: vec!["host-b.com".to_string()],
            destination_hosts: vec![WILDCARD_HOSTS.to_string()],
            port: 80,
            protocol: HTTP_PROTOCOL.to_string(),
        }];

        // Define restrictions for ns-04
        let ns_04_restrictions = vec![Restriction {
            hostnames: vec!["host-d.com".to_string()],
            destination_hosts: vec![WILDCARD_HOSTS.to_string()],
            port: 80,
            protocol: HTTP_PROTOCOL.to_string(),
        }];

        // Assign namespaces to gateways
        let gateway_01_namespaces: Namespaces = HashMap::from([
            (NS_01.to_string(), ns_01_restrictions),
            (NS_02.to_string(), ns_02_restrictions),
        ]);
        let gateway_02_namespaces: Namespaces =
            HashMap::from([(NS_04.to_string(), ns_04_restrictions)]);
        // Build the final restrictions map
        let restrictions = HashMap::from([
            (GW_01.to_string(), gateway_01_namespaces),
            (GW_02.to_string(), gateway_02_namespaces),
        ]);

        Settings {
            gateway_restrictions: restrictions,
        }
    }

    #[test]
    fn test_accept_allowed_host() -> Result<(), ()> {
        let settings = build_settings();
        let request_file = "test_data/vs-01.json";

        let tc = Testcase {
            name: String::from("accept allowed host"),
            fixture_file: String::from(request_file),
            expected_validation_result: true,
            settings,
        };

        let res = tc.eval(validate).unwrap();
        assert!(res.accepted);
        Ok(())
    }

    #[test]
    fn test_reject_disallowed_host() -> Result<(), ()> {
        let settings = build_settings();
        let request_file = "test_data/vs-02.json";

        let tc = Testcase {
            name: String::from("reject disallowed host"),
            fixture_file: String::from(request_file),
            expected_validation_result: false,
            settings,
        };
        let res = tc.eval(validate).unwrap();
        assert!(!res.accepted);
        Ok(())
    }

    #[test]
    fn test_reject_disallowed_services() -> Result<(), ()> {
        let ns_01_restrictions = vec![Restriction {
            hostnames: vec!["host-a.com".to_string()],
            destination_hosts: vec![
                "service-a.ns-01.svc.cluster.local".to_string(),
                "service-b.ns-01.svc.cluster.local".to_string(),
            ],
            port: 80,
            protocol: HTTP_PROTOCOL.to_string(),
        }];

        let gateway_02_namespaces: Namespaces =
            HashMap::from([(NS_01.to_string(), ns_01_restrictions)]);

        let restrictions = HashMap::from([("gateway-02".to_string(), gateway_02_namespaces)]);
        let settings = Settings {
            gateway_restrictions: restrictions,
        };

        let request_file = "test_data/vs-03.json";
        let tc = Testcase {
            name: String::from("reject disallowed destination host"),
            fixture_file: String::from(request_file),
            expected_validation_result: false,
            settings,
        };
        let res = tc.eval(validate).unwrap();
        assert!(!res.accepted);
        Ok(())
    }

    #[test]
    fn test_reject_disallowed_port() -> Result<(), ()> {
        let settings = build_settings();
        let request_file = "test_data/vs-04.json";

        let tc = Testcase {
            name: String::from("reject disallowed port"),
            fixture_file: String::from(request_file),
            expected_validation_result: false,
            settings,
        };
        let res = tc.eval(validate).unwrap();
        assert!(!res.accepted);
        Ok(())
    }

    #[test]
    fn test_reject_disallowed_tls() -> Result<(), ()> {
        let settings = build_settings();
        let request_file = "test_data/vs-05.json";

        let tc = Testcase {
            name: String::from("reject disallowed TLS protocol"),
            fixture_file: String::from(request_file),
            expected_validation_result: false,
            settings,
        };
        let res = tc.eval(validate).unwrap();
        assert!(!res.accepted);
        Ok(())
    }

    #[test]
    fn test_reject_disallowed_tcp() -> Result<(), ()> {
        let settings = build_settings();
        let request_file = "test_data/vs-06.json";

        let tc = Testcase {
            name: String::from("reject disallowed TCP route"),
            fixture_file: String::from(request_file),
            expected_validation_result: false,
            settings,
        };
        let res = tc.eval(validate).unwrap();
        assert!(!res.accepted);
        Ok(())
    }
}
