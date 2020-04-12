provider "k8s-yaml" {
  version = "1.0.0"
}

resource "k8s-yaml_raw" "test" {
  name = "dashboard"
  file_url = "https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0-rc7/aio/deploy/recommended.yaml"
}
