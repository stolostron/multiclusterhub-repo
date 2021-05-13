# Community Operator Package to Helm Chart Automation

This script takes in a set of parameters via a `config.yaml` file. The script parses these parameters in order to locate the referenced package.yaml file. Once located, the CSV of the operator is found, and parsed in order to create a helm chart.

## Configuration

All configuration for the script is contained in the `config.yaml` file. Each parameter is explained - 

```bash
- repo_name: "community-operators" # Name of the repository
  github_ref: "https://github.com/operator-framework/community-operators.git" # Repo to clone
  operators:
    - name: "assisted-service-operator" # Helm Chart Name
      channel: "ocm-2.3" # Channel (must exist in operator package.yaml file)
      package-yml: "community-operators/assisted-service-operator/assisted-service.package.yaml" # Link to the operators package.yaml, from base of repo
      imageMappings:
        assisted-service: assisted_service # Must contain a map to all images, even if the same. Otherwise script will fail
```

## Running

After the `config.yaml` file has been correctly configured, the script can be ran by calling `./csv-to-helm-chart.py` from the root directory the script is in. The associated helm chart will be created inside the `stable/` directory.

## Flags

### Destination

By adding the `destination` flag, the helm chart will be created at the specified directory. 

```bash
$ ./csv-to-helm-chart.py --destination .
```

### Skip Overrides

By adding the `skipOverrides` flag, overrides will not be applied, and instead will need to be applied manually. 

```bash
$ ./csv-to-helm-chart.py --skipOverrides=true
```