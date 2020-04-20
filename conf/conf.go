/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package conf

import (
	"flag"
	"fmt"
	"os"
)

type confResolver interface {
	Resolve(conf interface{}) error
}

func Load(p interface{}, yamlFlagName string) {

	initStruct(p)

	resolvers := make([]confResolver, 0)

	yamlFile := ""
	if yamlFlagName != "" {
		flag.StringVar(&yamlFile, yamlFlagName, "", "yaml config file")
	}

	er := &envResolver{}
	er.init(p)

	fr := &flagResolver{}
	fr.flagSet = flag.CommandLine
	fr.init(p)

	flag.Parse()

	if yamlFile != "" {
		yr := &yamlFileResolver{File: yamlFile}
		resolvers = append(resolvers, yr)
	}
	resolvers = append(resolvers, er)
	resolvers = append(resolvers, fr)

	for _, resolver := range resolvers {
		if err := resolver.Resolve(p); err != nil {
			fmt.Println(err.Error())
			os.Exit(2)
		}
	}

}
