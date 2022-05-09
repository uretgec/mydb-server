# Cobra Init
cd WORKING_DIR

### MainCommand
/Users/uretgec/go/bin/cobra add sitemap

### SubCommand
/Users/uretgec/go/bin/cobra add sitemapmydbserver -p "sitemapCmd" 

### RunCommand
go run main.go sitemap mydbserver --config .sitemap.dev.yaml

