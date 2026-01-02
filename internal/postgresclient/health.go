package postgresclient

import "context"

// HealthCheck 用于健康检查（k8s / lb）
func HealthCheck(ctx context.Context) error {
	sqlDB, err := DB().DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
