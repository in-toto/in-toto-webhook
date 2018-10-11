package pod_test 

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/slok/kubewebhook/pkg/log"
	hook "github.com/santiagotorres/in-toto-webhook/pkg/webhook"
)

func TestPodInTotovalidate(t *testing.T) {
    container := corev1.Container{
		Name: "test-container",
		Image: "namespace/working-test",
	}

	logger := &log.Std{}

	tests := []struct {
		name   string
		pod    *corev1.Pod
		shouldAdmit bool
		expErr bool
	}{
		{
			name: "Validating with proper metadata",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "test",
					Labels: map[string]string{"bruce": "wayne", "peter": "parker"},
				},
				Spec: corev1.PodSpec {
					Containers: []corev1.Container{container},
				},
			},
			shouldAdmit: false,
			expErr: false,
		},
	}
	validator := hook.PodValidator{
		Logger:   logger,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			gotPod := test.pod
			admit, _, err  := validator.Validate(context.TODO(), gotPod)

			assert.Equal(test.shouldAdmit, admit)
			if test.expErr {
				assert.Error(err)
			}
		})
	}
}
