//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package cloudstack

import (
	"encoding/json"
	"fmt"
	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"regexp"
	"strings"
)

func dataSourceCloudstackNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudstackNetworkRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			"can_use_for_deploy": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			// Computed values
			"network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCloudstackNetworkRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cloudstack.ListNetworksParams{}
	p.SetListall(true)
	p.SetCanusefordeploy(d.Get("can_use_for_deploy").(bool))

	csNetworks, err := cs.Network.ListNetworks(&p)
	if err != nil {
		return fmt.Errorf("failed to list networks: %s", err)
	}

	filters := d.Get("filter")
	var networks []*cloudstack.Network

	log.Printf("[DEBUG] Networks found: %d\n", len(csNetworks.Networks))

	for _, n := range csNetworks.Networks {
		log.Printf("[DEBUG] Checking network %s [id: %s]\n", n.Name, n.Id)
		match, err := applyNetworkFilters(n, filters.(*schema.Set))
		if err != nil {
			return err
		}

		if match {
			networks = append(networks, n)
		}
	}

	if len(networks) == 0 {
		return fmt.Errorf("no network is matching with the specified filters")
	}

	network := networks[0]

	log.Printf("[DEBUG] Selected network: %s\n", network.Displaytext)

	return networkDescriptionAttributes(d, network)
}

func networkDescriptionAttributes(d *schema.ResourceData, network *cloudstack.Network) error {
	d.SetId(network.Id)
	d.Set("network_id", network.Id)

	return nil
}

func applyNetworkFilters(network *cloudstack.Network, filters *schema.Set) (bool, error) {
	var networkJSON map[string]interface{}
	n, _ := json.Marshal(network)
	err := json.Unmarshal(n, &networkJSON)
	if err != nil {
		return false, err
	}

	for _, f := range filters.List() {
		m := f.(map[string]interface{})

		r, err := regexp.Compile(m["value"].(string))
		if err != nil {
			return false, fmt.Errorf("invalid regex: %s", err)
		}

		notFilter := strings.HasPrefix(m["name"].(string), "not_")
		fieldName := m["name"].(string)

		if notFilter {
			fieldName = strings.Replace(m["name"].(string), "not_", "", 1)
		}

		networkField := networkJSON[fieldName].(string)

		if notFilter && r.MatchString(networkField) {
			return false, nil
		}

		if !notFilter && !r.MatchString(networkField) {
			return false, nil
		}
	}

	return true, nil
}
