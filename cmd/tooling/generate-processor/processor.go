package main


var processorTemplateString string = "// package {{.PackageName}} is generated by generate-processor tooling\n// Make sure to insert real Description here\npackage {{.PackageName}}\n\n\nimport (\n    \"context\"\n    \"github.com/percybolmer/workflow/failure\"\n    \"github.com/percybolmer/workflow/metric\"\n    \"github.com/percybolmer/workflow/properties\"\n    \"github.com/percybolmer/workflow/relationships\"\n    \"github.com/percybolmer/workflow/processors\"\n\n)\n// {{.ProcessorName}} is used to $INSERT DESCRIPTION\ntype {{.ProcessorName}} struct{\n    Name     string\n    running  bool\n    cancel   context.CancelFunc\n    ingress  relationships.PayloadChannel\n    egress   relationships.PayloadChannel\n    failures relationships.FailurePipe\n    *properties.PropertyMap\n    *metric.Metrics\n}\n// init is really important since it is used to auto trigger Processor register when u import the generated package\nfunc init () {\n    err := processors.RegisterProcessor(\"{{.ProcessorName}}\", New{{.ProcessorName}}())\n    if err != nil {\n        panic(err)\n    }\n}\n\n// New{{.ProcessorName}} is used to initialize and generate a new processor\nfunc New{{.ProcessorName}}() *{{.ProcessorName}} {\n    proc := &{{.ProcessorName}}{\n        egress: make(relationships.PayloadChannel, 1000),\n        PropertyMap: properties.NewPropertyMap(),\n        Metrics: metric.NewMetrics(),\n    }\n\n    // Add Required Props -- remove_after\n    // proc.AddRequirement(\"your_property_names\")\n    return proc\n}\n\n// Initialize will make sure all needed Properties and Metrics are generated\nfunc (proc *{{.ProcessorName}}) Initialize() error {\n\n    // Make sure Properties are there\n    ok, _ := proc.ValidateProperties()\n    if !ok {\n        return properties.ErrRequiredPropertiesNotFulfilled\n    }\n    // If you need to read data from Properties and add to your Processor struct, this is the place to do it\n    return nil\n}\n\n// Start will spawn a goroutine that reads file and Exits either on Context.Done or When processing is finished\nfunc (proc *{{.ProcessorName}}) Start(ctx context.Context) error {\n    if proc.running {\n        return failure.ErrAlreadyRunning\n    }\n    // Uncomment if u need to Processor to require an Ingress relationship\n    //if rfp.ingress == nil {\n    //    return failure.ErrIngressRelationshipNeeded\n    //}\n\n    proc.running = true\n    // context will be used to spawn a Cancel func\n    c, cancel := context.WithCancel(ctx)\n    proc.cancel = cancel\n    go func() {\n        for {\n            select {\n                case payload := <-proc.ingress:\n                    // Do your processing here\n                    payload.GetPayload()\n\n                case <- c.Done():\n                    return\n            }\n        }\n    }()\n    return nil\n}\n\n// IsRunning will return true or false based on if the processor is currently running\nfunc (proc *{{.ProcessorName}}) IsRunning() bool {\n    return proc.running\n}\n// GetMetrics will return a bunch of generated metrics, or nil if there isn't any\nfunc (proc *{{.ProcessorName}}) GetMetrics() []*metric.Metric {\n    return proc.GetAllMetrics()\n}\n// SetFailureChannel will configure the failure channel of the Processor\nfunc (proc *{{.ProcessorName}}) SetFailureChannel(fp relationships.FailurePipe) {\n    proc.failures = fp\n}\n\n// Stop will stop the processing\nfunc (proc *{{.ProcessorName}}) Stop() {\n    if !proc.running {\n        return\n    }\n    proc.running = false\n    proc.cancel()\n}\n// SetIngress will change the ingress of the processor, Restart is needed before applied changes\nfunc (proc *{{.ProcessorName}}) SetIngress(i relationships.PayloadChannel) {\n    proc.ingress = i\n    return\n}\n// GetEgress will return an Egress that is used to output the processed items\nfunc (proc *{{.ProcessorName}}) GetEgress() relationships.PayloadChannel {\n    return proc.egress\n}"


var testTemplateString string = "package {{.PackageName}}\n\nimport (\n    \"testing\"\n    \"context\"\n)\n\nfunc Test{{.ProcessorName}}_Initialize(t *testing.T) {\n    proc := New{{.ProcessorName}}()\n    // Empty Init should fail\n    proc.AddRequirement(\"test_prop\")\n    err := proc.Initialize()\n    if err == nil {\n        t.Fatal(\"Should have failed to initialize without proper Properties set\")\n    }\n    proc.SetProperty(\"test_prop\", true)\n    err = proc.Initialize()\n    if err != nil {\n        t.Fatalf(\"Should have passed initialize now that all properties are there: %s\", err.Error())\n    }\n\n}\n\nfunc Test{{.ProcessorName}}_StopStart(t *testing.T){\n    proc := New{{.ProcessorName}}()\n\n\n    if proc.IsRunning(){\n        t.Fatal(\"proccessor should not be running after creation\")\n    }\n\n    err := proc.Start(context.TODO())\n    if err != nil {\n        t.Fatal(\"processor should not have failed to startup\")\n    }\n    err = proc.Start(context.TODO())\n    if !errors.Is(err, failure.ErrAlreadyRunning) {\n        t.Fatal(\"processor should have reported already running\")\n    }\n\n    proc.Stop()\n    err = proc.Start(context.TODO())\n    if err != nil {\n        t.Fatal(\"a restart should have been able to run\")\n    }\n}"