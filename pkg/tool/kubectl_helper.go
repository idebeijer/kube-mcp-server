package tool

import (
	"context"
	"os/exec"
)

func (h *Handler) runKubectl(ctx context.Context, args ...string) ([]byte, error) {
	if h.kubeconfigPath != "" {
		args = append([]string{"--kubeconfig", h.kubeconfigPath}, args...)
	}

	cmd := exec.CommandContext(ctx, h.kubectlPath, args...)
	return cmd.CombinedOutput()
}
