# Changes

Code differences compared to source project.

## cmd/demo2kratos/main.go (+11 -4)

```diff
@@ -4,15 +4,16 @@
 	"flag"
 	"log/slog"
 	"os"
+	"path/filepath"
 
 	"github.com/go-kratos/kratos/contrib/otel/v3/tracing"
 	"github.com/go-kratos/kratos/v3"
 	"github.com/go-kratos/kratos/v3/config"
-	"github.com/go-kratos/kratos/v3/config/file"
 	"github.com/go-kratos/kratos/v3/log"
 	"github.com/go-kratos/kratos/v3/transport/grpc"
 	"github.com/go-kratos/kratos/v3/transport/http"
 	"github.com/yylego/done"
+	"github.com/yylego/kratos-config/configkratos"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/conf"
 	"github.com/yylego/must"
 	"github.com/yylego/rese"
@@ -62,10 +63,16 @@
 		slog.String("service.version", Version),
 	)
 	log.SetDefault(logger)
+	// demo2 uses DataSource to load dynamic config with watch support
+	// demo2 使用 DataSource 加载动态配置，支持 Watch 机制
+	var sources []config.Source
+	for _, item := range rese.A1(os.ReadDir(flagconf)) {
+		configPath := filepath.Join(flagconf, item.Name())
+		configData := rese.A1(os.ReadFile(configPath))
+		sources = append(sources, configkratos.NewDataSource(configData, "yaml"))
+	}
 	c := config.New(
-		config.WithSource(
-			file.NewSource(flagconf),
-		),
+		config.WithSource(sources...),
 	)
 	defer rese.F0(c.Close)
 
```

