go build -o terraform-provider-k8s-yaml.exe

Copy-Item ".\terraform-provider-k8s-yaml.exe" -Destination "${env:APPDATA}\\terraform.d\\plugins\\terraform-provider-k8s-yaml_v1.0.0.exe" -Force
