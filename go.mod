module github.com/rancher/opni-monitoring

go 1.18

require (
	emperror.dev/errors v0.8.1
	github.com/Masterminds/sprig/v3 v3.2.2
	github.com/alecthomas/jsonschema v0.0.0-20220216202328-9eeeec9d044b
	github.com/andybalholm/brotli v1.0.4
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de
	github.com/banzaicloud/operator-tools v0.28.4
	github.com/cert-manager/cert-manager v1.8.0
	github.com/cortexproject/cortex v1.12.0-rc.0
	github.com/dghubble/trie v0.0.0-20211002190126-ca25329b35c6
	github.com/go-logr/logr v1.2.3
	github.com/gofiber/fiber/v2 v2.31.0
	github.com/golang/mock v1.6.0
	github.com/google/go-cmp v0.5.7
	github.com/google/uuid v1.3.0
	github.com/hashicorp/go-hclog v1.2.0
	github.com/hashicorp/go-plugin v1.4.3
	github.com/iancoleman/strcase v0.2.0
	github.com/jedib0t/go-pretty/v6 v6.3.0
	github.com/jhump/protoreflect v1.12.0
	github.com/jwalton/go-supportscolor v1.1.0
	github.com/kralicky/gpkg v0.0.0-20220311205216-0d8ea9557555
	github.com/kralicky/grpc-gateway/v2 v2.7.3-0.20220201000610-57444701bbdc
	github.com/kralicky/kmatch v0.0.0-20210910033132-e5a80a7a45e6
	github.com/kralicky/spellbook v0.0.0-20220322004259-9900ddc1e90a
	github.com/lestrrat-go/backoff/v2 v2.0.8
	github.com/lestrrat-go/jwx v1.2.22
	github.com/magefile/mage v1.13.0
	github.com/mattn/go-tty v0.0.4
	github.com/mitchellh/mapstructure v1.4.3
	github.com/olebedev/when v0.0.0-20211212231525-59bd4edcf9d6
	github.com/onsi/ginkgo/v2 v2.1.3
	github.com/onsi/gomega v1.19.0
	github.com/ory/fosite v0.42.1
	github.com/phayes/freeport v0.0.0-20220201140144-74d24b5ae9f5
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.55.1
	github.com/prometheus/client_golang v1.12.1
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.33.0
	github.com/prometheus/prometheus v1.8.2-0.20211119115433-692a54649ed7
	github.com/rancher/opni v0.4.0
	github.com/samber/lo v1.12.0
	github.com/schollz/progressbar/v3 v3.8.6
	github.com/spf13/cobra v1.4.0
	github.com/ttacon/chalk v0.0.0-20160626202418-22c06c80ed31
	github.com/valyala/fasthttp v1.35.0
	github.com/vearutop/statigz v1.1.8
	go.etcd.io/etcd/client/v3 v3.5.2
	go.etcd.io/etcd/etcdctl/v3 v3.5.2
	go.uber.org/atomic v1.9.0
	go.uber.org/zap v1.21.0
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4
	golang.org/x/exp v0.0.0-20220407100705-7b9b53b0aca4
	golang.org/x/mod v0.6.0-dev.0.20211013180041-c96bc1413d57
	gonum.org/v1/gonum v0.11.0
	google.golang.org/genproto v0.0.0-20220329172620-7be39ac1afc7
	google.golang.org/grpc v1.45.0
	google.golang.org/protobuf v1.28.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api v0.23.5
	k8s.io/apiextensions-apiserver v0.23.5
	k8s.io/apimachinery v0.23.5
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/utils v0.0.0-20220210201930-3a6ce19ff2f9
	opensearch.opster.io v0.0.0-00010101000000-000000000000
	sigs.k8s.io/controller-runtime v0.11.2
	sigs.k8s.io/controller-tools v0.8.0
	sigs.k8s.io/yaml v1.3.0
)

require (
	github.com/AlekSi/pointer v1.0.0 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.1.1 // indirect
	github.com/NVIDIA/gpu-operator v1.8.1 // indirect
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/aws/aws-sdk-go v1.43.10 // indirect
	github.com/banzaicloud/k8s-objectmatcher v1.8.0 // indirect
	github.com/banzaicloud/logging-operator/pkg/sdk v0.7.19 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bgentry/speakeasy v0.1.0 // indirect
	github.com/briandowns/spinner v1.12.0 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/cppforlife/go-patch v0.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.0-20210816181553-5444fa50b93d // indirect
	github.com/dennwc/varint v1.0.0 // indirect
	github.com/dgraph-io/ristretto v0.0.3 // indirect
	github.com/dsnet/compress v0.0.2-0.20210315054119-f66993602bf5 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/edsrzf/mmap-go v1.1.0 // indirect
	github.com/evanphx/json-patch v5.6.0+incompatible // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/form3tech-oss/jwt-go v3.2.3+incompatible // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-kit/kit v0.12.0 // indirect
	github.com/go-kit/log v0.2.0 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0 // indirect
	github.com/gobuffalo/flect v0.2.3 // indirect
	github.com/goccy/go-json v0.9.6 // indirect
	github.com/gogo/googleapis v1.4.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/gogo/status v1.1.0 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/pprof v0.0.0-20220218203455-0368bd9e19a7 // indirect
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/yamux v0.0.0-20180604194846-3520598351bb // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/iancoleman/orderedmap v0.2.0 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.15.0 // indirect
	github.com/klauspost/pgzip v1.2.5 // indirect
	github.com/kralicky/ragu v0.2.7-0.20220323195938-430b56cdb12d // indirect
	github.com/lestrrat-go/blackmagic v1.0.0 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/iter v1.0.1 // indirect
	github.com/lestrrat-go/option v1.0.0 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mattn/goveralls v0.0.6 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/mholt/archiver/v4 v4.0.0-alpha.3 // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-testing-interface v1.0.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/nwaples/rardecode/v2 v2.0.0-beta.2 // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/opentracing-contrib/go-grpc v0.0.0-20210225150812-73cb765af46e // indirect
	github.com/opentracing-contrib/go-stdlib v1.0.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/ory/go-acc v0.2.6 // indirect
	github.com/ory/go-convenience v0.1.0 // indirect
	github.com/ory/viper v1.7.5 // indirect
	github.com/ory/x v0.0.214 // indirect
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/pierrec/lz4/v4 v4.1.12 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/common/sigv4 v0.1.0 // indirect
	github.com/prometheus/node_exporter v1.0.0-rc.0.0.20200428091818-01054558c289 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/sercand/kuberesolver v2.4.0+incompatible // indirect
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.7.1 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/therootcompany/xz v1.0.1 // indirect
	github.com/uber/jaeger-client-go v2.29.1+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	github.com/wayneashleyberry/terminal-dimensions v1.0.0 // indirect
	github.com/weaveworks/common v0.0.0-20210913144402-035033b78a78 // indirect
	github.com/weaveworks/promrus v1.2.0 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	github.com/yoheimuta/go-protoparser/v4 v4.5.4 // indirect
	go.etcd.io/bbolt v1.3.6 // indirect
	go.etcd.io/etcd/api/v3 v3.5.2 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.2 // indirect
	go.etcd.io/etcd/client/v2 v2.305.2 // indirect
	go.etcd.io/etcd/etcdutl/v3 v3.5.2 // indirect
	go.etcd.io/etcd/pkg/v3 v3.5.2 // indirect
	go.etcd.io/etcd/raft/v3 v3.5.2 // indirect
	go.etcd.io/etcd/server/v3 v3.5.2 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.28.0 // indirect
	go.opentelemetry.io/otel v1.4.1 // indirect
	go.opentelemetry.io/otel/trace v1.4.1 // indirect
	go.uber.org/goleak v1.1.12 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f // indirect
	golang.org/x/oauth2 v0.0.0-20220223155221-ee480838109b // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20220330033206-e17cdc41300f // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20220210224613-90d013bbcef8 // indirect
	golang.org/x/tools v0.1.9 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gomodules.xyz/jsonpatch/v2 v2.2.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/cheggaaa/pb.v1 v1.0.28 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.66.2 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/component-base v0.23.5 // indirect
	k8s.io/klog/v2 v2.40.1 // indirect
	k8s.io/kube-openapi v0.0.0-20220114203427-a0453230fd26 // indirect
	sigs.k8s.io/gateway-api v0.4.1 // indirect
	sigs.k8s.io/json v0.0.0-20211208200746-9f7c6b3444d2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.1 // indirect
)

replace (
	github.com/NVIDIA/gpu-operator => github.com/kralicky/gpu-operator v1.8.1-0.20211112183255-72529edf38be

	github.com/banzaicloud/logging-operator/pkg/sdk => github.com/banzaicloud/logging-operator/pkg/sdk v0.0.0-20220225205714-b06e7ad17676

	github.com/openshift/api => github.com/openshift/api v0.0.0-20210216211028-bb81baaf35cd
	go.uber.org/zap => github.com/kralicky/zap v1.19.2-0.20220311060549-c0d473b28cca
	// Because of a dependency chain to Cortex
	k8s.io/client-go => k8s.io/client-go v0.23.5
	opensearch.opster.io => github.com/dbason/opensearch-k8s-operator/opensearch-operator v0.0.0-20220405232819-82b73065a2ca
)
