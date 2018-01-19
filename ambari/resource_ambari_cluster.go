package ambari

import (
	"fmt"
	restClient "github.com/disaster37/go-ambari-rest"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func resourceAmbariCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAmbariClusterCreate,
		Read:   resourceAmbariClusterRead,
		Update: resourceAmbariClusterUpdate,
		Delete: resourceAmbariClusterDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAmbariClusterImport,
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAmbariClusterCreate(d *schema.ResourceData, meta interface{}) error {
	log.Info("Creating Cluster")
	clusterName := d.Get("cluster_name").(string)
	version := d.Get("version").(string)

	client := meta.(*Client).Client()

	cluster := &restClient.Cluster{
		Cluster: &restClient.ClusterInfo{
			Version:     version,
			ClusterName: clusterName,
		},
	}
	log.Debugf("Cluster to create: %s", cluster.String())

	cluster, err := client.CreateCluster(cluster)
	if err != nil {
		return err
	}
	log.Debugf("Cluster after to create: %s", cluster.String())

	// Use cluster name as id, because of ambari manage cluster like it
	d.SetId(cluster.Cluster.ClusterName)
	log.Infof("Cluster ID: %s", d.Id())

	return resourceAmbariClusterRead(d, meta)
}

func resourceAmbariClusterRead(d *schema.ResourceData, meta interface{}) error {
	log.Infof("Refreshing Cluster: %s", d.Id())

	clusterName := d.Id()
	client := meta.(*Client).Client()
	log.Debugf("ClusterName: %s", clusterName)

	cluster, err := client.Cluster(clusterName)
	if err != nil {
		return err
	}

	if cluster == nil {
		return errors.New(fmt.Sprintf("Cluster with name %s not found", clusterName))
	}

	log.Debug("Cluster after read: %s", cluster.String())
	log.Infof("Cluster ID : %s", cluster.Cluster.ClusterName)

	d.Set("cluster_name", cluster.Cluster.ClusterName)
	d.Set("cluster_id", strconv.FormatInt(cluster.Cluster.ClusterId, 10))
	d.Set("version", cluster.Cluster.Version)

	return nil
}

func resourceAmbariClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Infof("Updating Cluster: %s", d.Id())
	oldClusterName := d.Id()
	clusterName := d.Get("cluster_name").(string)
	client := meta.(*Client).Client()

	cluster := &restClient.Cluster{
		Cluster: &restClient.ClusterInfo{
			ClusterName: clusterName,
		},
	}
	log.Debugf("Cluster to update: %s", cluster.String())

	cluster, err := client.UpdateCluster(oldClusterName, cluster)
	if err != nil {
		return err
	}
	log.Debugf("Cluster after update: %s", cluster.String())
	log.Infof("New Cluster ID: %s", cluster.Cluster.ClusterName)

	d.SetId(cluster.Cluster.ClusterName)

	return resourceAmbariClusterRead(d, meta)
}

func resourceAmbariClusterDelete(d *schema.ResourceData, meta interface{}) error {
	log.Infof("Deleting Privilege: %s", d.Id())
	clusterName := d.Id()
	client := meta.(*Client).Client()

	log.Debugf("ClusterName: %s", clusterName)
	err := client.DeleteCluster(clusterName)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceAmbariClusterImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	err := resourceAmbariClusterRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
