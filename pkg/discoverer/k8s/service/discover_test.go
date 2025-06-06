//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package service manages the main logic of server.
package service

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		selector *config.Selectors
// 		opts     []Option
// 	}
// 	type want struct {
// 		wantDsc Discoverer
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Discoverer, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotDsc Discoverer, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotDsc, w.wantDsc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotDsc, w.wantDsc)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           selector:nil,
// 		           opts:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           selector:nil,
// 		           opts:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			gotDsc, err := New(test.args.selector, test.args.opts...)
// 			if err := checkFunc(test.want, gotDsc, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_discoverer_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		maxPods         int
// 		nodes           sync.Map[string, *node.Node]
// 		nodeMetrics     sync.Map[string, mnode.Node]
// 		pods            sync.Map[string, *[]pod.Pod]
// 		podMetrics      sync.Map[string, mpod.Pod]
// 		services        sync.Map[string, *service.Service]
// 		podsByNode      atomic.Pointer[map[string]map[string]map[string][]*payload.Info_Pod]
// 		podsByNamespace atomic.Pointer[map[string]map[string][]*payload.Info_Pod]
// 		podsByName      atomic.Pointer[map[string][]*payload.Info_Pod]
// 		nodeByName      atomic.Pointer[map[string]*payload.Info_Node]
// 		svcsByName      atomic.Pointer[map[string]*payload.Info_Service]
// 		ctrl            k8s.Controller
// 		namespace       string
// 		name            string
// 		csd             time.Duration
// 		der             net.Dialer
// 		eg              errgroup.Group
// 	}
// 	type want struct {
// 		want <-chan error
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, <-chan error, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got <-chan error, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           maxPods:0,
// 		           nodes:nil,
// 		           nodeMetrics:nil,
// 		           pods:nil,
// 		           podMetrics:nil,
// 		           services:nil,
// 		           podsByNode:nil,
// 		           podsByNamespace:nil,
// 		           podsByName:nil,
// 		           nodeByName:nil,
// 		           svcsByName:nil,
// 		           ctrl:nil,
// 		           namespace:"",
// 		           name:"",
// 		           csd:nil,
// 		           der:nil,
// 		           eg:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           maxPods:0,
// 		           nodes:nil,
// 		           nodeMetrics:nil,
// 		           pods:nil,
// 		           podMetrics:nil,
// 		           services:nil,
// 		           podsByNode:nil,
// 		           podsByNamespace:nil,
// 		           podsByName:nil,
// 		           nodeByName:nil,
// 		           svcsByName:nil,
// 		           ctrl:nil,
// 		           namespace:"",
// 		           name:"",
// 		           csd:nil,
// 		           der:nil,
// 		           eg:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			d := &discoverer{
// 				maxPods:         test.fields.maxPods,
// 				nodes:           test.fields.nodes,
// 				nodeMetrics:     test.fields.nodeMetrics,
// 				pods:            test.fields.pods,
// 				podMetrics:      test.fields.podMetrics,
// 				services:        test.fields.services,
// 				podsByNode:      test.fields.podsByNode,
// 				podsByNamespace: test.fields.podsByNamespace,
// 				podsByName:      test.fields.podsByName,
// 				nodeByName:      test.fields.nodeByName,
// 				svcsByName:      test.fields.svcsByName,
// 				ctrl:            test.fields.ctrl,
// 				namespace:       test.fields.namespace,
// 				name:            test.fields.name,
// 				csd:             test.fields.csd,
// 				der:             test.fields.der,
// 				eg:              test.fields.eg,
// 			}
//
// 			got, err := d.Start(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_discoverer_GetPods(t *testing.T) {
// 	type args struct {
// 		req *payload.Discoverer_Request
// 	}
// 	type fields struct {
// 		maxPods         int
// 		nodes           sync.Map[string, *node.Node]
// 		nodeMetrics     sync.Map[string, mnode.Node]
// 		pods            sync.Map[string, *[]pod.Pod]
// 		podMetrics      sync.Map[string, mpod.Pod]
// 		services        sync.Map[string, *service.Service]
// 		podsByNode      atomic.Pointer[map[string]map[string]map[string][]*payload.Info_Pod]
// 		podsByNamespace atomic.Pointer[map[string]map[string][]*payload.Info_Pod]
// 		podsByName      atomic.Pointer[map[string][]*payload.Info_Pod]
// 		nodeByName      atomic.Pointer[map[string]*payload.Info_Node]
// 		svcsByName      atomic.Pointer[map[string]*payload.Info_Service]
// 		ctrl            k8s.Controller
// 		namespace       string
// 		name            string
// 		csd             time.Duration
// 		der             net.Dialer
// 		eg              errgroup.Group
// 	}
// 	type want struct {
// 		wantPods *payload.Info_Pods
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Info_Pods, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotPods *payload.Info_Pods, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotPods, w.wantPods) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotPods, w.wantPods)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           maxPods:0,
// 		           nodes:nil,
// 		           nodeMetrics:nil,
// 		           pods:nil,
// 		           podMetrics:nil,
// 		           services:nil,
// 		           podsByNode:nil,
// 		           podsByNamespace:nil,
// 		           podsByName:nil,
// 		           nodeByName:nil,
// 		           svcsByName:nil,
// 		           ctrl:nil,
// 		           namespace:"",
// 		           name:"",
// 		           csd:nil,
// 		           der:nil,
// 		           eg:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           maxPods:0,
// 		           nodes:nil,
// 		           nodeMetrics:nil,
// 		           pods:nil,
// 		           podMetrics:nil,
// 		           services:nil,
// 		           podsByNode:nil,
// 		           podsByNamespace:nil,
// 		           podsByName:nil,
// 		           nodeByName:nil,
// 		           svcsByName:nil,
// 		           ctrl:nil,
// 		           namespace:"",
// 		           name:"",
// 		           csd:nil,
// 		           der:nil,
// 		           eg:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			d := &discoverer{
// 				maxPods:         test.fields.maxPods,
// 				nodes:           test.fields.nodes,
// 				nodeMetrics:     test.fields.nodeMetrics,
// 				pods:            test.fields.pods,
// 				podMetrics:      test.fields.podMetrics,
// 				services:        test.fields.services,
// 				podsByNode:      test.fields.podsByNode,
// 				podsByNamespace: test.fields.podsByNamespace,
// 				podsByName:      test.fields.podsByName,
// 				nodeByName:      test.fields.nodeByName,
// 				svcsByName:      test.fields.svcsByName,
// 				ctrl:            test.fields.ctrl,
// 				namespace:       test.fields.namespace,
// 				name:            test.fields.name,
// 				csd:             test.fields.csd,
// 				der:             test.fields.der,
// 				eg:              test.fields.eg,
// 			}
//
// 			gotPods, err := d.GetPods(test.args.req)
// 			if err := checkFunc(test.want, gotPods, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_discoverer_GetNodes(t *testing.T) {
// 	type args struct {
// 		req *payload.Discoverer_Request
// 	}
// 	type fields struct {
// 		maxPods         int
// 		nodes           sync.Map[string, *node.Node]
// 		nodeMetrics     sync.Map[string, mnode.Node]
// 		pods            sync.Map[string, *[]pod.Pod]
// 		podMetrics      sync.Map[string, mpod.Pod]
// 		services        sync.Map[string, *service.Service]
// 		podsByNode      atomic.Pointer[map[string]map[string]map[string][]*payload.Info_Pod]
// 		podsByNamespace atomic.Pointer[map[string]map[string][]*payload.Info_Pod]
// 		podsByName      atomic.Pointer[map[string][]*payload.Info_Pod]
// 		nodeByName      atomic.Pointer[map[string]*payload.Info_Node]
// 		svcsByName      atomic.Pointer[map[string]*payload.Info_Service]
// 		ctrl            k8s.Controller
// 		namespace       string
// 		name            string
// 		csd             time.Duration
// 		der             net.Dialer
// 		eg              errgroup.Group
// 	}
// 	type want struct {
// 		wantNodes *payload.Info_Nodes
// 		err       error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Info_Nodes, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotNodes *payload.Info_Nodes, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotNodes, w.wantNodes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotNodes, w.wantNodes)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           maxPods:0,
// 		           nodes:nil,
// 		           nodeMetrics:nil,
// 		           pods:nil,
// 		           podMetrics:nil,
// 		           services:nil,
// 		           podsByNode:nil,
// 		           podsByNamespace:nil,
// 		           podsByName:nil,
// 		           nodeByName:nil,
// 		           svcsByName:nil,
// 		           ctrl:nil,
// 		           namespace:"",
// 		           name:"",
// 		           csd:nil,
// 		           der:nil,
// 		           eg:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           maxPods:0,
// 		           nodes:nil,
// 		           nodeMetrics:nil,
// 		           pods:nil,
// 		           podMetrics:nil,
// 		           services:nil,
// 		           podsByNode:nil,
// 		           podsByNamespace:nil,
// 		           podsByName:nil,
// 		           nodeByName:nil,
// 		           svcsByName:nil,
// 		           ctrl:nil,
// 		           namespace:"",
// 		           name:"",
// 		           csd:nil,
// 		           der:nil,
// 		           eg:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			d := &discoverer{
// 				maxPods:         test.fields.maxPods,
// 				nodes:           test.fields.nodes,
// 				nodeMetrics:     test.fields.nodeMetrics,
// 				pods:            test.fields.pods,
// 				podMetrics:      test.fields.podMetrics,
// 				services:        test.fields.services,
// 				podsByNode:      test.fields.podsByNode,
// 				podsByNamespace: test.fields.podsByNamespace,
// 				podsByName:      test.fields.podsByName,
// 				nodeByName:      test.fields.nodeByName,
// 				svcsByName:      test.fields.svcsByName,
// 				ctrl:            test.fields.ctrl,
// 				namespace:       test.fields.namespace,
// 				name:            test.fields.name,
// 				csd:             test.fields.csd,
// 				der:             test.fields.der,
// 				eg:              test.fields.eg,
// 			}
//
// 			gotNodes, err := d.GetNodes(test.args.req)
// 			if err := checkFunc(test.want, gotNodes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_discoverer_GetServices(t *testing.T) {
// 	type args struct {
// 		req *payload.Discoverer_Request
// 	}
// 	type fields struct {
// 		maxPods         int
// 		nodes           sync.Map[string, *node.Node]
// 		nodeMetrics     sync.Map[string, mnode.Node]
// 		pods            sync.Map[string, *[]pod.Pod]
// 		podMetrics      sync.Map[string, mpod.Pod]
// 		services        sync.Map[string, *service.Service]
// 		podsByNode      atomic.Pointer[map[string]map[string]map[string][]*payload.Info_Pod]
// 		podsByNamespace atomic.Pointer[map[string]map[string][]*payload.Info_Pod]
// 		podsByName      atomic.Pointer[map[string][]*payload.Info_Pod]
// 		nodeByName      atomic.Pointer[map[string]*payload.Info_Node]
// 		svcsByName      atomic.Pointer[map[string]*payload.Info_Service]
// 		ctrl            k8s.Controller
// 		namespace       string
// 		name            string
// 		csd             time.Duration
// 		der             net.Dialer
// 		eg              errgroup.Group
// 	}
// 	type want struct {
// 		wantSvcs *payload.Info_Services
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Info_Services, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotSvcs *payload.Info_Services, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotSvcs, w.wantSvcs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotSvcs, w.wantSvcs)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           maxPods:0,
// 		           nodes:nil,
// 		           nodeMetrics:nil,
// 		           pods:nil,
// 		           podMetrics:nil,
// 		           services:nil,
// 		           podsByNode:nil,
// 		           podsByNamespace:nil,
// 		           podsByName:nil,
// 		           nodeByName:nil,
// 		           svcsByName:nil,
// 		           ctrl:nil,
// 		           namespace:"",
// 		           name:"",
// 		           csd:nil,
// 		           der:nil,
// 		           eg:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           maxPods:0,
// 		           nodes:nil,
// 		           nodeMetrics:nil,
// 		           pods:nil,
// 		           podMetrics:nil,
// 		           services:nil,
// 		           podsByNode:nil,
// 		           podsByNamespace:nil,
// 		           podsByName:nil,
// 		           nodeByName:nil,
// 		           svcsByName:nil,
// 		           ctrl:nil,
// 		           namespace:"",
// 		           name:"",
// 		           csd:nil,
// 		           der:nil,
// 		           eg:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			d := &discoverer{
// 				maxPods:         test.fields.maxPods,
// 				nodes:           test.fields.nodes,
// 				nodeMetrics:     test.fields.nodeMetrics,
// 				pods:            test.fields.pods,
// 				podMetrics:      test.fields.podMetrics,
// 				services:        test.fields.services,
// 				podsByNode:      test.fields.podsByNode,
// 				podsByNamespace: test.fields.podsByNamespace,
// 				podsByName:      test.fields.podsByName,
// 				nodeByName:      test.fields.nodeByName,
// 				svcsByName:      test.fields.svcsByName,
// 				ctrl:            test.fields.ctrl,
// 				namespace:       test.fields.namespace,
// 				name:            test.fields.name,
// 				csd:             test.fields.csd,
// 				der:             test.fields.der,
// 				eg:              test.fields.eg,
// 			}
//
// 			gotSvcs, err := d.GetServices(test.args.req)
// 			if err := checkFunc(test.want, gotSvcs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
