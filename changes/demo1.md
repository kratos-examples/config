# Changes

Code differences compared to source project.

## cmd/demo1kratos/main.go (+12 -4)

```diff
@@ -3,15 +3,16 @@
 import (
 	"flag"
 	"os"
+	"path/filepath"
 
 	"github.com/go-kratos/kratos/v2"
 	"github.com/go-kratos/kratos/v2/config"
-	"github.com/go-kratos/kratos/v2/config/file"
 	"github.com/go-kratos/kratos/v2/log"
 	"github.com/go-kratos/kratos/v2/middleware/tracing"
 	"github.com/go-kratos/kratos/v2/transport/grpc"
 	"github.com/go-kratos/kratos/v2/transport/http"
 	"github.com/yylego/done"
+	"github.com/yylego/kratos-config/configkratos"
 	"github.com/yylego/kratos-examples/demo1kratos/internal/conf"
 	"github.com/yylego/must"
 	"github.com/yylego/rese"
@@ -56,10 +57,17 @@
 		"trace.id", tracing.TraceID(),
 		"span.id", tracing.SpanID(),
 	)
+
+	// demo1 uses DataStatic to load static config
+	// demo1 使用 DataStatic 加载静态配置
+	var sources []config.Source
+	for _, item := range rese.A1(os.ReadDir(flagconf)) {
+		configPath := filepath.Join(flagconf, item.Name())
+		configData := rese.A1(os.ReadFile(configPath))
+		sources = append(sources, configkratos.NewDataStatic(configData, "yaml"))
+	}
 	c := config.New(
-		config.WithSource(
-			file.NewSource(flagconf),
-		),
+		config.WithSource(sources...),
 	)
 	defer rese.F0(c.Close)
 
```

