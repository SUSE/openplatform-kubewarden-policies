use crate::LOG_DRAIN;
use serde::{Deserialize, Serialize};
use slog::info;
use std::collections::HashMap;

/// Represents the top-level JSON structure.
#[derive(Serialize, Deserialize, Default, Debug)]
#[serde(default, rename_all = "camelCase")]
pub struct Settings {
    pub gateway_restrictions: HashMap<Gateway, Namespaces>,
}

/// A type alias for clarity, representing a gateway from a string
pub type Gateway = String;

/// A type alias for clarity, representing a map from a namespace string
/// to a vector of `Restriction` rules.
pub type Namespaces = HashMap<String, Vec<Restriction>>;

/// Represents a single restriction rule with hostnames and destination hosts.
#[derive(Serialize, Deserialize, Debug)]
pub struct Restriction {
    /// A list of hostnames for the restriction.
    /// The `alias` attribute allows serde to deserialize from "hostname" as well.
    /// The `default` attribute handles cases where the field might be missing.
    #[serde(alias = "hostname", default)]
    pub hostnames: Vec<String>,

    /// A list of destination hosts for the restriction.
    /// For additional details on the naming convention please check:
    /// https://istio.io/latest/docs/reference/config/networking/virtual-service/#Headers
    #[serde(default)]
    pub destination_hosts: Vec<String>,

    /// The destination port for the restriction
    #[serde(default)]
    pub port: i32,

    /// The destination protocol for the restriction
    #[serde(default)]
    pub protocol: String,
}

impl kubewarden::settings::Validatable for Settings {
    fn validate(&self) -> Result<(), String> {
        info!(LOG_DRAIN, "starting settings validation");
        // Settings bindings cannot be empty
        if self.gateway_restrictions.is_empty() {
            return Err("The 'gatewayRestrictions' map cannot be empty.".to_string());
        }
        // Iterate over each element and check for empty fields.
        for (gateway, namespaces) in &self.gateway_restrictions {
            if gateway.is_empty() {
                return Err(
                    "Gateway names inside 'gatewayRestrictions' cannot be empty.".to_string(),
                );
            }
            if namespaces.is_empty() {
                return Err(format!(
                    "The namespace list for gateway '{gateway}' cannot be empty.",
                ));
            }
        }

        info!(LOG_DRAIN, "settings validation successful");
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    static HTTP_PROTOCOL: &str = "http";
    static GW_01: &str = "gateway-01";
    static NS_01: &str = "ns-01";

    use kubewarden_policy_sdk::settings::Validatable;

    /// Helper function to create a valid Restriction for tests
    fn create_valid_restriction() -> Restriction {
        Restriction {
            hostnames: vec!["host.example.com".to_string()],
            destination_hosts: vec!["service-a".to_string()],
            port: 80,
            protocol: HTTP_PROTOCOL.to_string(),
        }
    }

    #[test]
    fn test_valid_settings() {
        let settings = Settings {
            gateway_restrictions: HashMap::from([(
                GW_01.to_string(),
                HashMap::from([(NS_01.to_string(), vec![create_valid_restriction()])]),
            )]),
        };
        assert!(settings.validate().is_ok());
    }

    #[test]
    fn test_empty_gateway_restrictions_is_invalid() {
        let settings = Settings {
            gateway_restrictions: HashMap::new(),
        };
        let result = settings.validate();
        assert!(result.is_err());
        assert_eq!(
            result.unwrap_err(),
            "The 'gatewayRestrictions' map cannot be empty."
        );
    }

    #[test]
    fn test_empty_gateway_name_is_invalid() {
        let settings = Settings {
            gateway_restrictions: HashMap::from([(
                "".to_string(), // Empty gateway name
                HashMap::from([(NS_01.to_string(), vec![create_valid_restriction()])]),
            )]),
        };
        let result = settings.validate();
        assert!(result.is_err());
        assert_eq!(
            result.unwrap_err(),
            "Gateway names inside 'gatewayRestrictions' cannot be empty."
        );
    }

    #[test]
    fn test_empty_namespace_list_is_invalid() {
        let settings = Settings {
            gateway_restrictions: HashMap::from([(
                GW_01.to_string(),
                HashMap::new(), // Empty namespace map
            )]),
        };
        let result = settings.validate();
        assert!(result.is_err());
        assert_eq!(
            result.unwrap_err(),
            "The namespace list for gateway 'gateway-01' cannot be empty."
        );
    }
}
