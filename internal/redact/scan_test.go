package redact

import (
	"reflect"
	"testing"

	"github.com/smol-utils/actiondoc/internal/model"
)

// wfWithRun wraps a single run: script in a minimal workflow.
func wfWithRun(run string) *model.Workflow {
	return &model.Workflow{Jobs: []model.Job{{ID: "j", Steps: []model.Step{{Run: run}}}}}
}

// scanHosts runs walkHostsURLs over s and returns what it classified.
func scanHosts(s string) (urls, hosts []string) {
	walkHostsURLs(s,
		func(u string) { urls = append(urls, u) },
		func(h string) { hosts = append(hosts, h) })
	return
}

// TestWalkHostsURLs_Precision locks the bare-hostname classifier against the dotted
// tokens that crowd shell-heavy run: scripts. The false-positive rows are real cases
// observed on the dependency-track dogfood repo: Java -D system properties, Maven
// coordinates, archive filenames, and git config keys all match the hostname shape but
// must not be redacted. Over-redaction is safe but makes redacted output unreadable.
func TestWalkHostsURLs_Precision(t *testing.T) {
	cases := []struct {
		name  string
		in    string
		hosts []string
		urls  []string
	}{
		// Genuine hosts must still be caught.
		{"public host", "git log --format=%ae shows users.noreply.github.com", []string{"users.noreply.github.com"}, nil},
		{"registry host", "docker push docker.io/foo/bar", []string{"docker.io"}, nil},
		{"internal pseudo-TLD", "scp build.tgz deploy.internal.corp:/srv", []string{"deploy.internal.corp"}, nil},
		{"cluster-local", "curl backend.default.svc.cluster.local", []string{"backend.default.svc.cluster.local"}, nil},
		{"scheme URL unaffected by TLD rule", "curl https://nexus.corp.mycompany/repo", nil, []string{"https://nexus.corp.mycompany/repo"}},

		// dependency-track false positives.
		{"java -D property", "java -Dlogback.configuration.file=logback.xml -jar app.jar", nil, nil},
		{"maven coordinate", "mvn org.codehaus.mojo:exec-maven-plugin:exec", nil, nil},
		{"bundled jar", "cp dependency-track-bundled.jar /opt", nil, nil},
		{"dockerfile variant", "docker build -f Dockerfile.alpine .", nil, nil},
		{"git config keys", "git config user.email x; git config user.name y", nil, nil},

		// Adjacent shapes that must stay non-hosts.
		{"lowercase property no TLD", "set logback.configuration.file please", nil, nil},
		{"version string", "upgrade to 1.2.3 now", nil, nil},
		{"filename", "edit config.yaml then run deploy.sh", nil, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urls, hosts := scanHosts(tc.in)
			if !reflect.DeepEqual(hosts, tc.hosts) {
				t.Errorf("hosts: got %v want %v", hosts, tc.hosts)
			}
			if !reflect.DeepEqual(urls, tc.urls) {
				t.Errorf("urls: got %v want %v", urls, tc.urls)
			}
		})
	}
}

// TestApply_HostPrecisionEndToEnd confirms the classifier is applied symmetrically:
// a run: script mixing false-positive shapes with a genuine host comes out with only
// the host replaced.
func TestApply_HostPrecisionEndToEnd(t *testing.T) {
	in := "java -Dlogback.configuration.file=x -jar dependency-track-bundled.jar && " +
		"git config user.email ci@x && curl registry.internal.corp/health"
	w := wfWithRun(in)
	Apply(wf(w), Options{})
	got := w.Jobs[0].Steps[0].Run
	want := "java -Dlogback.configuration.file=x -jar dependency-track-bundled.jar && " +
		"git config user.email ci@x && curl HOST_1/health"
	if got != want {
		t.Errorf("got  %q\nwant %q", got, want)
	}
}
