package gnmit

import (
	"context"
	"time"

	log "github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

const (
	// Unix socket for central datastore.
	DatastoreAddress = "/tmp/datastore.api"
)

type DatastoreServer struct {
	gpb.UnimplementedGNMIServer // For forward-compatibility

	gnmiServer *GNMIServer
}

func NewDatastoreServer(gnmiServer *GNMIServer) *DatastoreServer {
	return &DatastoreServer{gnmiServer: gnmiServer}
}

func (d *DatastoreServer) Set(_ context.Context, req *gpb.SetRequest) (*gpb.SetResponse, error) {
	// TODO(wenbli): Reject values that modify config values. We only allow modifying state through this server.
	// TODO(wenbli): Unmarshal to a struct without PreferShadowPath in
	// order to validate that state paths are valid according to the
	// schema.
	if err := d.gnmiServer.set(req, false); err != nil {
		return &gpb.SetResponse{}, status.Errorf(codes.Aborted, "%v", err)
	}

	// SetRequest has been validated, so we update the cache.
	deletes := append([]*gpb.Path{}, req.Delete...)
	for _, update := range req.Replace {
		deletes = append(deletes, update.Path)
	}
	t := d.gnmiServer.c.cache.GetTarget(d.gnmiServer.c.name)
	notif := &gpb.Notification{
		Timestamp: time.Now().UnixNano(),
		Prefix:    req.Prefix,
		Delete:    deletes,
		Update:    req.Update,
	}
	if notif.Prefix.Origin == "" {
		notif.Prefix.Origin = OpenconfigOrigin
	}
	log.V(1).Infof("datastore updates central cache: %v", notif)
	t.GnmiUpdate(notif)
	return &gpb.SetResponse{}, nil
}