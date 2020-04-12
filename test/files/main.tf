provider "k8s-yaml" {
  version = "1.0.0"
}

resource "k8s-yaml_raw" "test" {
  name = "dashboard"
  files = [
    "${path.module}/service.yaml"
  ]
}
