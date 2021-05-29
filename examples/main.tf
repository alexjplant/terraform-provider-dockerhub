terraform {
  required_providers {
    dockerhub = {
      source = "github.com/alexjplant/dockerhub"
    }
  }
}

provider "dockerhub" {
  
}

data "dockerhub_image_tag" "nginx_latest" {
  repository_name = "nginx"
  namespace = "library"
  tag_name = "latest"
}
