package valid

// +kubebuilder4helm:rbac:groups=batch.io,resources=cronjobs,verbs=get;watch;create
// +kubebuilder4helm:rbac:groups=batch.io,resources=cronjobs/status,verbs=get;update;patch
// +kubebuilder4helm:rbac:groups=art,resources=jobs,verbs=get
// +kubebuilder4helm:rbac:groups=wave,resources=jobs,verbs=get,namespace=zoo
// +kubebuilder4helm:rbac:groups=batch;batch;batch,resources=jobs/status,verbs=watch
// +kubebuilder4helm:rbac:groups=batch;cron,resources=jobs/status,verbs=create;get
// +kubebuilder4helm:rbac:groups=art,resources=jobs,verbs=get,namespace=zoo
// +kubebuilder4helm:rbac:groups=cron;batch,resources=jobs/status,verbs=get;create
// +kubebuilder4helm:rbac:groups=batch,resources=jobs/status,verbs=watch;watch
// +kubebuilder4helm:rbac:groups=art,resources=jobs,verbs=get,namespace=park
// +kubebuilder4helm:rbac:groups=batch.io,resources=cronjobs,resourceNames=foo;bar;baz,verbs=get;watch
