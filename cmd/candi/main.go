package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

const (
	ps1         = "\x1b[32;1m>>> \x1b[0m"
	packageName = "pkg.agungdwiprasetyo.com/candi"
	initService = "initservice"
	addModule   = "addModule"
)

var (
	scopeMap = map[string]string{
		"1": "initservice", "2": "addmodule",
	}
	serviceHandlersMap = map[string]string{
		"1": "restHandler", "2": "grpcHandler", "3": "graphqlHandler",
	}
	workerHandlersMap = map[string]string{
		"1": "kafkaHandler", "2": "schedulerHandler", "3": "redissubsHandler", "4": "taskqueueHandler",
	}
	tpl *template.Template
)

type param struct {
	PackageName string
	ServiceName string
	Modules     []string
}

// FileStructure model
type FileStructure struct {
	TargetDir    string
	IsDir        bool
	FromTemplate bool
	DataSource   interface{}
	Source       string
	FileName     string
	Skip         bool
	Childs       []FileStructure
}

func main() {

	scope, serviceName, modules, serviceHandlers, workerHandlers := parseInput()

	var data param
	data.PackageName = packageName
	data.ServiceName = serviceName

	dataSourceWithHandler := map[string]string{"PackageName": packageName, "ServiceName": serviceName}
	mergeMap(dataSourceWithHandler, serviceHandlers)
	mergeMap(dataSourceWithHandler, workerHandlers)

	tpl = template.New(packageName)

	apiStructure := FileStructure{
		TargetDir: "api/", IsDir: true,
	}

	cmdStructure := FileStructure{
		TargetDir: "cmd/", IsDir: true,
		Childs: []FileStructure{
			{
				TargetDir: "{{.ServiceName}}/", IsDir: true, DataSource: data,
				Childs: []FileStructure{
					{FromTemplate: true, DataSource: data, Source: cmdMainTemplate, FileName: "main.go"},
					{FromTemplate: true, DataSource: dataSourceWithHandler, Source: envTemplate, FileName: ".env"},
					{FromTemplate: true, DataSource: dataSourceWithHandler, Source: envTemplate, FileName: ".env.sample"},
				},
			},
		},
	}

	serviceStructure := FileStructure{
		TargetDir: "internal/", IsDir: true, DataSource: data,
	}

	apiProtoStructure := FileStructure{
		TargetDir: "proto/", IsDir: true,
	}
	apiGraphQLStructure := FileStructure{
		TargetDir: "graphql/", IsDir: true,
	}

	if scope == addModule {
		files, err := ioutil.ReadDir("internal/modules")
		if err != nil {
			panic(err)
		}
		for _, f := range files {
			if f.IsDir() {
				data.Modules = append(data.Modules, f.Name())
			}
		}
	}

	var moduleStructure = FileStructure{
		TargetDir: "modules/", IsDir: true, DataSource: data,
	}
	for _, moduleName := range modules {
		moduleName = strings.TrimSpace(moduleName)
		data.Modules = append(data.Modules, moduleName)
		dataSource := map[string]string{"module": moduleName}
		mergeMap(dataSource, dataSourceWithHandler)
		fmt.Println("map {}", dataSource)

		cleanArchModuleDir := []FileStructure{
			{
				TargetDir: "delivery/", IsDir: true,
				Childs: []FileStructure{
					{TargetDir: "graphqlhandler/", IsDir: true, Childs: []FileStructure{
						{FromTemplate: true, DataSource: dataSource, Source: deliveryGraphqlRootTemplate, FileName: "root_resolver.go"},
						{FromTemplate: true, DataSource: dataSource, Source: deliveryGraphqlQueryTemplate, FileName: "query_resolver.go"},
						{FromTemplate: true, DataSource: dataSource, Source: deliveryGraphqlMutationTemplate, FileName: "mutation_resolver.go"},
						{FromTemplate: true, DataSource: dataSource, Source: deliveryGraphqlSubscriptionTemplate, FileName: "subscription_resolver.go"},
					}},
					{TargetDir: "grpchandler/", IsDir: true, Childs: []FileStructure{
						{FromTemplate: true, DataSource: dataSource, Source: deliveryGRPCTemplate, FileName: "grpchandler.go"},
					}},
					{TargetDir: "resthandler/", IsDir: true, Childs: []FileStructure{
						{FromTemplate: true, DataSource: dataSource, Source: deliveryRestTemplate, FileName: "resthandler.go"},
					}},
					{TargetDir: "workerhandler/", IsDir: true, Childs: []FileStructure{
						{FromTemplate: true, DataSource: dataSource, Source: deliveryKafkaTemplate, FileName: "kafka_handler.go"},
						{FromTemplate: true, DataSource: dataSource, Source: deliveryRedisTemplate, FileName: "redis_handler.go"},
						{FromTemplate: true, DataSource: dataSource, Source: deliveryCronTemplate, FileName: "cron_handler.go"},
						{FromTemplate: true, DataSource: dataSource, Source: deliveryTaskQueueTemplate, FileName: "taskqueue_handler.go"},
					}},
				},
			},
			{
				TargetDir: "domain/", IsDir: true,
				Childs: []FileStructure{
					{FromTemplate: true, FileName: "domain.go"},
				},
			},
			{
				TargetDir: "repository/", IsDir: true,
				Childs: []FileStructure{
					{TargetDir: "interfaces/", IsDir: true},
					{FromTemplate: true, FileName: "repository.go"},
				},
			},
			{
				TargetDir: "usecase/", IsDir: true,
				Childs: []FileStructure{
					{FromTemplate: true, FileName: "usecase.go"},
					{FromTemplate: true, FileName: "usecase_impl.go"},
				},
			},
		}

		moduleStructure.Childs = append(moduleStructure.Childs, []FileStructure{
			{
				TargetDir: moduleName + "/", IsDir: true,
				Childs: append(cleanArchModuleDir,
					FileStructure{
						FromTemplate: true, DataSource: dataSource, Source: moduleMainTemplate, FileName: "module.go",
					},
				),
			},
		}...)

		apiProtoStructure.Childs = append(apiProtoStructure.Childs, FileStructure{
			TargetDir: moduleName, IsDir: true,
			Childs: []FileStructure{
				{FromTemplate: true, DataSource: dataSource, Source: defaultGRPCProto, FileName: moduleName + ".proto"},
			},
		})
		apiGraphQLStructure.Childs = append(apiGraphQLStructure.Childs, FileStructure{
			FromTemplate: true, DataSource: dataSource, Source: defaultGraphqlSchema, FileName: moduleName + ".graphql",
		})
	}
	serviceStructure.Childs = append(serviceStructure.Childs, moduleStructure)
	serviceStructure.Childs = append(serviceStructure.Childs, FileStructure{
		FromTemplate: true, DataSource: data, Source: serviceMainTemplate, FileName: "service.go",
	})

	apiGraphQLStructure.Childs = append(apiGraphQLStructure.Childs, FileStructure{
		FromTemplate: true, DataSource: data, Source: defaultGraphqlRootSchema, FileName: "_schema.graphql",
	})
	apiStructure.Childs = []FileStructure{
		apiGraphQLStructure,
		{
			TargetDir: "jsonschema/", IsDir: true,
			Childs: []FileStructure{
				// {FromTemplate: true, DataSource: dataSource, FileName: "_schema.json"},
			},
		},
		apiProtoStructure,
	}

	var baseDirectoryFile FileStructure
	baseDirectoryFile.TargetDir = "{{.ServiceName}}/"
	baseDirectoryFile.DataSource = data
	baseDirectoryFile.IsDir = true
	switch scope {
	case initService:
		baseDirectoryFile.Childs = []FileStructure{
			apiStructure, cmdStructure, serviceStructure,
			{FromTemplate: true, DataSource: data, Source: dockerfileTemplate, FileName: "Dockerfile"},
			{FromTemplate: true, DataSource: data, Source: makefileTemplate, FileName: "Makefile"},
			{FromTemplate: true, DataSource: data, Source: gomodTemplate, FileName: "go.mod"},
			{Source: gitignoreTemplate, FileName: ".gitignore"},
		}

	case addModule:
		moduleStructure.Skip = true
		serviceStructure.Skip = true
		serviceStructure.Childs = []FileStructure{
			moduleStructure,
			{FromTemplate: true, DataSource: data, Source: serviceMainTemplate, FileName: "service.go"},
		}

		apiStructure.Skip = true
		apiProtoStructure.Skip, apiGraphQLStructure.Skip = true, true
		apiStructure.Childs = []FileStructure{
			apiProtoStructure, apiGraphQLStructure,
		}

		baseDirectoryFile.Childs = []FileStructure{apiStructure, serviceStructure}
		baseDirectoryFile.Skip = true

	default:
		panic("invalid scope parameter")
	}

	exec(baseDirectoryFile)
}

func exec(fl FileStructure) {
	dirBuff := loadTemplate(fl.TargetDir, fl.DataSource)
	dirName := string(dirBuff)

	if fl.Skip {
		goto execChild
	}

	if _, err := os.Stat(dirName); os.IsExist(err) {
		panic(err)
	}

	if fl.IsDir {
		fmt.Printf("creating %s...\n", dirName)
		if err := os.Mkdir(dirName, 0700); err != nil {
			fmt.Println("mkdir err:", err)
			panic(err)
		}
	}

	if fl.FileName != "" {
		var buff []byte
		if fl.FromTemplate {
			if fl.Source != "" {
				buff = loadTemplate(fl.Source, fl.DataSource)
			} else {
				lastDir := filepath.Dir(fl.TargetDir)
				buff = defaultDataSource(lastDir[strings.LastIndex(lastDir, "/")+1:])
			}
		} else {
			buff = []byte(fl.Source)
		}
		dirName = strings.TrimSuffix(dirName, "/")
		if err := ioutil.WriteFile(dirName+"/"+fl.FileName, buff, 0644); err != nil {
			panic(err)
		}
	}

execChild:
	for _, child := range fl.Childs {
		child.TargetDir = dirName + child.TargetDir
		exec(child)
	}
}

func loadTemplate(source string, sourceData interface{}) []byte {
	var byteBuff = new(bytes.Buffer)
	defer byteBuff.Reset()

	tmpl, err := tpl.Funcs(formatTemplate()).Parse(source)
	if err != nil {
		panic(err)
	}

	if err := tmpl.Execute(byteBuff, sourceData); err != nil {
		panic(err)
	}

	return byteBuff.Bytes()
}

func formatTemplate() template.FuncMap {
	replacer := strings.NewReplacer("-", "", "*", "", "/", "", ":", "")
	return template.FuncMap{

		"clean": func(v string) string {
			return replacer.Replace(v)
		},

		"upper": func(str string) string {
			return strings.Title(str)
		},

		"isHandlerActive": func(str string) string {
			ok, _ := strconv.ParseBool(str)
			if ok {
				return ""
			}
			return "// "
		},
	}
}

func parseInput() (scope, serviceName string, modules []string, serviceHandlers, workerHandlers map[string]string) {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\033[1mWhat do you want?\n" +
		"1) Init service\n" +
		"2) Add module(s)\033[0m")

	cmdInput, _ := reader.ReadString('\n')
	cmdInput = strings.TrimRight(cmdInput, "\n")
	var ok bool
	scope, ok = scopeMap[cmdInput]
	if !ok {
		fmt.Println("invalid option")
		os.Exit(1)
	}

	if scope == initService {
		fmt.Print(ps1 + "\033[1mPlease input service name:\033[0m ")
		cmdInput, _ := reader.ReadString('\n')
		cmdInput = strings.TrimRight(cmdInput, "\n")
		serviceName = cmdInput

	}

	fmt.Print(ps1 + "\033[1mPlease input module names (separated by comma):\033[0m ")
	cmdInput, _ = reader.ReadString('\n')
	cmdInput = strings.TrimRight(cmdInput, "\n")
	modules = strings.Split(cmdInput, ",")
	sort.Slice(modules, func(i, j int) bool {
		return modules[i] < modules[j]
	})

	fmt.Print(ps1 + "\033[1mPlease select service handlers (separated by comma)\n" +
		"0) all\n" +
		"1) Rest API\n" +
		"2) GRPC\n" +
		"3) GraphQL\033[0m\n")
	cmdInput, _ = reader.ReadString('\n')
	cmdInput = strings.TrimRight(cmdInput, "\n")

	serviceHandlers = make(map[string]string, 3)
	for i := 1; i <= 3; i++ {
		serviceHandlers[serviceHandlersMap[strconv.Itoa(i)]] = "false"
	}
	for _, str := range strings.Split(strings.Trim(cmdInput, ","), ",") {
		serviceHandlers[serviceHandlersMap[strings.TrimSpace(str)]] = "true"
	}

	fmt.Print(ps1 + "\033[1mPlease select worker handlers (separated by comma)\n" +
		"0) all\n" +
		"1) Kafka Consumer\n" +
		"2) Scheduler\n" +
		"3) Redis Subscriber\n" +
		"4) Task Queue\033[0m\n")
	cmdInput, _ = reader.ReadString('\n')
	cmdInput = strings.TrimRight(cmdInput, "\n")

	workerHandlers = make(map[string]string, 4)
	for i := 1; i <= 4; i++ {
		workerHandlers[workerHandlersMap[strconv.Itoa(i)]] = "false"
	}
	for _, str := range strings.Split(strings.Trim(cmdInput, ","), ",") {
		workerHandlers[workerHandlersMap[strings.TrimSpace(str)]] = "true"
	}

	return
}

func mergeMap(dest, source map[string]string) {
	for k, v := range source {
		dest[k] = v
	}
}