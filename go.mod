module github.com/davidcollom/terratest-utils

go 1.26.3

// Keep gateway-api pinned for compatibility with the Linkerd controller APIs
// used by pkg/linkerd. Newer Linkerd releases use non-semver tags, so the
// github.com/linkerd/linkerd2 module is pinned explicitly below.
replace sigs.k8s.io/gateway-api => sigs.k8s.io/gateway-api v1.0.0

replace github.com/argoproj/argo-cd/gitops-engine => github.com/argoproj/argo-cd/gitops-engine v0.0.0-20260508195653-5ec06031ba40

require (
	github.com/argoproj/argo-cd/gitops-engine v0.7.1-0.20250908182407-97ad5b59a627
	github.com/argoproj/argo-cd/v3 v3.4.5
	github.com/argoproj/argo-events v1.9.11
	github.com/argoproj/argo-rollouts v1.9.1
	github.com/argoproj/argo-workflows/v3 v3.7.17
	github.com/cert-manager/cert-manager v1.20.2
	github.com/external-secrets/external-secrets/apis v0.0.0-20260407212151-e325bced502e
	github.com/fluxcd/helm-controller/api v1.5.3
	github.com/fluxcd/kustomize-controller/api v1.8.3
	github.com/fluxcd/source-controller/api v1.8.2
	github.com/gruntwork-io/terratest v0.56.0
	github.com/linkerd/linkerd2 v0.5.1-0.20260622225159-eadc1acf79ad
	github.com/stretchr/testify v1.11.1
	github.com/vmware-tanzu/velero v1.18.2
	istio.io/api v1.29.1
	istio.io/client-go v1.29.1
	k8s.io/api v0.35.4
	k8s.io/apiextensions-apiserver v0.35.4
	k8s.io/apimachinery v0.35.4
	k8s.io/client-go v0.35.4
	sigs.k8s.io/controller-runtime v0.23.3
)

require (
	cloud.google.com/go/compute/metadata v0.9.0 // indirect
	cyphar.com/go-pathrs v0.2.4 // indirect
	dario.cat/mergo v1.0.2 // indirect
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.21.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.13.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.12.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.7.1 // indirect
	github.com/MakeNowJust/heredoc v1.0.0 // indirect
	github.com/Masterminds/semver/v3 v3.5.0 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/ProtonMail/go-crypto v1.4.1 // indirect
	github.com/argoproj/pkg v0.13.7-0.20250305113207-cbc37dc61de5 // indirect
	github.com/argoproj/pkg/v2 v2.0.1 // indirect
	github.com/aws/aws-sdk-go-v2 v1.41.12 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.10 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.32.17 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.19.16 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.23 // indirect
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.22.18 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.23 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.23 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.24 // indirect
	github.com/aws/aws-sdk-go-v2/service/acm v1.38.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/autoscaling v1.65.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.68.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.57.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.297.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ecr v1.57.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ecs v1.77.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/iam v1.53.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.11.21 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.23 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.23 // indirect
	github.com/aws/aws-sdk-go-v2/service/kms v1.50.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/lambda v1.89.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/rds v1.117.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/route53 v1.62.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.101.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.41.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.0.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/sns v1.39.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/sqs v1.42.25 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssm v1.68.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.21 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.42.1 // indirect
	github.com/aws/smithy-go v1.27.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/bmatcuk/doublestar/v4 v4.10.0 // indirect
	github.com/bombsimon/logrusr/v4 v4.1.0 // indirect
	github.com/boombuler/barcode v1.1.0 // indirect
	github.com/bradleyfalzon/ghinstallation/v2 v2.18.0 // indirect
	github.com/casbin/casbin/v2 v2.135.0 // indirect
	github.com/casbin/govaluate v1.10.0 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/chai2010/gettext-go v1.0.3 // indirect
	github.com/chainguard-dev/git-urls v1.0.2 // indirect
	github.com/cloudflare/circl v1.6.3 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/cyphar/filepath-securejoin v0.6.1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/dlclark/regexp2 v1.11.5 // indirect
	github.com/emicklei/go-restful/v3 v3.13.0 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/evanphx/json-patch/v5 v5.9.11 // indirect
	github.com/exponent-io/jsonpath v0.0.0-20210407135951-1de76d718b3f // indirect
	github.com/fatih/camelcase v1.0.0 // indirect
	github.com/fluxcd/pkg/apis/acl v0.9.0 // indirect
	github.com/fluxcd/pkg/apis/kustomize v1.16.0 // indirect
	github.com/fluxcd/pkg/apis/meta v1.26.0 // indirect
	github.com/fsnotify/fsnotify v1.10.1 // indirect
	github.com/fxamacker/cbor/v2 v2.9.1 // indirect
	github.com/go-errors/errors v1.5.1 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-git/go-billy/v5 v5.9.0 // indirect
	github.com/go-git/go-git/v5 v5.19.1 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-openapi/jsonpointer v0.23.1 // indirect
	github.com/go-openapi/jsonreference v0.21.6 // indirect
	github.com/go-openapi/swag v0.26.1 // indirect
	github.com/go-openapi/swag/cmdutils v0.26.1 // indirect
	github.com/go-openapi/swag/conv v0.26.1 // indirect
	github.com/go-openapi/swag/fileutils v0.26.1 // indirect
	github.com/go-openapi/swag/jsonname v0.26.1 // indirect
	github.com/go-openapi/swag/jsonutils v0.26.1 // indirect
	github.com/go-openapi/swag/loading v0.26.1 // indirect
	github.com/go-openapi/swag/mangling v0.26.1 // indirect
	github.com/go-openapi/swag/netutils v0.26.1 // indirect
	github.com/go-openapi/swag/stringutils v0.26.1 // indirect
	github.com/go-openapi/swag/typeutils v0.26.1 // indirect
	github.com/go-openapi/swag/yamlutils v0.26.1 // indirect
	github.com/go-redis/cache/v9 v9.0.0 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/gobwas/glob v0.2.4-0.20181002190808-e7a84e9525fe // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/btree v1.1.3 // indirect
	github.com/google/gnostic-models v0.7.1 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/go-github/v69 v69.2.0 // indirect
	github.com/google/go-github/v84 v84.0.0 // indirect
	github.com/google/go-querystring v1.2.0 // indirect
	github.com/google/uuid v1.6.1-0.20241114170450-2d3c2a9cc518 // indirect
	github.com/gorilla/websocket v1.5.4-0.20250319132907-e064f32e3674 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/gruntwork-io/go-commons v0.17.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.9.2 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jonboulle/clockwork v0.5.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/kevinburke/ssh_config v1.6.0 // indirect
	github.com/klauspost/compress v1.18.6 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de // indirect
	github.com/mattn/go-zglob v0.0.6 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/moby/spdystream v0.5.1 // indirect
	github.com/moby/term v0.5.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.3-0.20250322232337-35a7c28c31ee // indirect
	github.com/monochromegane/go-gitignore v0.0.0-20200626010858-205db1a8cc00 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mxk/go-flowrate v0.0.0-20140419014527-cca7078d478f // indirect
	github.com/nxadm/tail v1.4.11 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/patrickmn/go-cache v2.1.1-0.20191004192108-46f407853014+incompatible // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pjbgf/sha1cd v0.6.0 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/pquerna/otp v1.5.0 // indirect
	github.com/prometheus/client_golang v1.23.2 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.68.1 // indirect
	github.com/prometheus/procfs v0.20.1 // indirect
	github.com/redis/go-redis/v9 v9.18.0 // indirect
	github.com/robfig/cron/v3 v3.0.2-0.20210106135023-bc59245fe10e // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sergi/go-diff v1.4.0 // indirect
	github.com/sirupsen/logrus v1.9.4 // indirect
	github.com/skeema/knownhosts v1.3.2 // indirect
	github.com/spf13/cobra v1.10.2 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/urfave/cli/v2 v2.27.7 // indirect
	github.com/vmihailenco/go-tinylfu v0.2.2 // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	github.com/xlab/treeprint v1.2.0 // indirect
	github.com/xrash/smetrics v0.0.0-20250705151800-55b8f293f342 // indirect
	go.opentelemetry.io/otel v1.44.0 // indirect
	go.opentelemetry.io/otel/trace v1.44.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.yaml.in/yaml/v2 v2.4.4 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/crypto v0.54.0 // indirect
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/oauth2 v0.36.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/term v0.45.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	golang.org/x/time v0.15.0 // indirect
	gomodules.xyz/jsonpatch/v2 v2.5.0 // indirect
	google.golang.org/genproto v0.0.0-20260406210006-6f92a3bedf2d // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260615183401-62b3387ff324 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260630182238-925bb5da69e7 // indirect
	google.golang.org/grpc v1.82.1 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/evanphx/json-patch.v4 v4.13.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apiserver v0.35.4 // indirect
	k8s.io/cli-runtime v0.35.3 // indirect
	k8s.io/component-base v0.35.4 // indirect
	k8s.io/component-helpers v0.35.3 // indirect
	k8s.io/controller-manager v0.35.3 // indirect
	k8s.io/klog/v2 v2.140.0 // indirect
	k8s.io/kube-aggregator v0.35.4 // indirect
	k8s.io/kube-openapi v0.0.0-20260330154417-16be699c7b31 // indirect
	k8s.io/kubectl v0.35.3 // indirect
	k8s.io/kubernetes v1.35.3 // indirect
	k8s.io/utils v0.0.0-20260319190234-28399d86e0b5 // indirect
	oras.land/oras-go/v2 v2.6.2 // indirect
	sigs.k8s.io/gateway-api v1.5.1 // indirect
	sigs.k8s.io/json v0.0.0-20250730193827-2d320260d730 // indirect
	sigs.k8s.io/kustomize/api v0.21.1 // indirect
	sigs.k8s.io/kustomize/kyaml v0.21.1 // indirect
	sigs.k8s.io/randfill v1.0.0 // indirect
	sigs.k8s.io/structured-merge-diff/v6 v6.3.2 // indirect
	sigs.k8s.io/yaml v1.6.0 // indirect
)
