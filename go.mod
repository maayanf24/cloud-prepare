module github.com/submariner-io/cloud-prepare

go 1.13

require (
	github.com/aws/aws-sdk-go v1.38.51
	github.com/golang/mock v1.4.4
	github.com/submariner-io/admiral v0.10.0-m1
	golang.org/x/oauth2 v0.0.0-20210427180440-81ed05c6b58c
	google.golang.org/api v0.30.0
	k8s.io/apimachinery v0.21.0
	k8s.io/client-go v1.5.2
)

// Pinned to kubernetes-1.19.10
replace (
	k8s.io/api => k8s.io/api v0.19.10
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.10
	k8s.io/client-go => k8s.io/client-go v0.19.10
)
