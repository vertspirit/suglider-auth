package rbac

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/casbin/casbin/v2"
	"github.com/jmoiron/sqlx"
	"github.com/memwey/casbin-sqlx-adapter"
)

type CasbinSettings struct {
	Config      string
	Table       string
	Db          *sqlx.DB
	EnableCache bool
}

type CasbinEnforcerConfig struct {
	Enforcer    *casbin.CachedEnforcer
	CasbinTable string
}

type CasbinPolicy struct {
	Sub         string    `json:"subject"`
	Obj         string    `json:"object"`
	Act         string    `json:"action"`
}

type CasbinGroupingPolicy struct {
	Member      string    `json:"member"`
	Role        string    `json:"role"`
}

type CasbinObject struct {
	Obj         string    `json:"object"`
}

func NewCasbinCachedEnforcer(cs *CasbinSettings) (*casbin.CachedEnforcer, error) {
	csbnAdapterOpts := &sqlxadapter.AdapterOptions {
		DB:        cs.Db,
		TableName: cs.Table,
		// DriverName:     "mysql",
		// DataSourceName: "root:1234@tcp(127.0.0.1:3306)/suglider",
	}
	csbnAdapter := sqlxadapter.NewAdapterFromOptions(csbnAdapterOpts)
	csbnEnforcer, err := casbin.NewCachedEnforcer(cs.Config, csbnAdapter)
	if err != nil {
		return nil, err
	}
	return csbnEnforcer, nil
}

func NewCasbinEnforcerConfig(cs *CasbinSettings) (*CasbinEnforcerConfig, error) {
	enforcer, err := NewCasbinCachedEnforcer(cs)
	if err != nil {
		return nil, err
	}
	enforcer.EnableAutoSave(true)
	enforcer.EnableCache(cs.EnableCache)
	csbnConfig := &CasbinEnforcerConfig {
		Enforcer:    enforcer,
		CasbinTable: cs.Table,
	}
	return csbnConfig, nil
}

func(cec *CasbinEnforcerConfig) InitPolicies() error {
	if ok, err := cec.Enforcer.Enforcer.AddPolicy("admin", "/*", "*"); !ok {
		if err != nil {
			return err
		}
		slog.Info("This policy already exists.")
	}
	anonymousPolicies := []string{
		"/static",
		"/login",
		"/sign-up",
		"/api/v1/user/login",
		"/api/v1/user/logout",
		"/api/v1/user/sign-up",
	}
	for _, item := range anonymousPolicies {
		if ok, err := cec.Enforcer.Enforcer.AddPolicy("anonymous", item, "GET"); !ok {
			if err != nil {
				return err
			}
			slog.Info("This policy already exists.")
		}
		if ok, err := cec.Enforcer.Enforcer.AddPolicy("anonymous", item, "POST"); !ok {
			if err != nil {
				return err
			}
			slog.Info("This policy already exists.")
		}
	}
	cec.Enforcer.LoadPolicy()
	return nil
}

func(cs *CasbinSettings) ListPoliciesCtx(ctx context.Context) ([][]string, error) {
	policies := make([][]string, 0)
	query := fmt.Sprintf("SELECT DISTINCT v0, v1, v2 FROM %s WHERE %s = ?", cs.Table, "p_type")
	rows, err := cs.Db.QueryContext(ctx, query, "p")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var (
			sub string
			obj string
			act string
		)
		if err := rows.Scan(&sub, &obj, &act); err != nil {
			return nil, err
		}
		policies = append(policies, []string { sub, obj, act })
	}
	return policies, nil
}

func(cs *CasbinSettings) ListGroupingPoliciesCtx(ctx context.Context) ([][]string, error) {
	gpolicies := make([][]string, 0)
	query := fmt.Sprintf("SELECT DISTINCT v0, v1 FROM %s WHERE %s = ?", cs.Table, "p_type")
	rows, err := cs.Db.QueryContext(ctx, query, "g")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var (
			m  string
			r  string
		)
		if err := rows.Scan(&m, &r); err != nil {
			return nil, err
		}
		gpolicies = append(gpolicies, []string { m, r })
	}
	return gpolicies, nil
}

func(cs *CasbinSettings) ListRolesCtx(ctx context.Context) ([]string, error) {
	roles := make([]string, 0)
	query := fmt.Sprintf("SELECT DISTINCT %s FROM %s WHERE %s = ?", "*", cs.Table, "p_type")
	rows, err := cs.Db.QueryContext(ctx, query, "g")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return roles, nil
}

func(cs *CasbinSettings) ListMembersCtx(ctx context.Context) ([]string, error) {
	members := make([]string, 0)
	query := fmt.Sprintf("SELECT DISTINCT %s FROM %s WHERE %s = ?", "v0", cs.Table, "p_type")
	rows, err := cs.Db.QueryContext(ctx, query, "g")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		members = append(members, role)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return members, nil
}

func(cs *CasbinSettings) GetMembersWithRoleCtx(ctx context.Context, name string) ([]string, error) {
	members := make([]string, 0)
	query := fmt.Sprintf(
		"SELECT DISTINCT %s FROM %s WHERE %s = ? AND %s = ?",
		"v0",
		cs.Table,
		"p_type",
		"v1",
	)
	rows, err := cs.Db.QueryContext(ctx, query, "g", name)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var account string
		if err := rows.Scan(&account); err != nil {
			return nil, err
		}
		members = append(members, account)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return members, nil
}

func(cs *CasbinSettings) GetRolesOfMemberCtx(ctx context.Context, name string) ([]string, error) {
	roles := make([]string, 0)
	query := fmt.Sprintf(
		"SELECT DISTINCT %s FROM %s WHERE %s = ? AND %s = ?",
		"v1",
		cs.Table,
		"p_type",
		"v0",
	)
	rows, err := cs.Db.QueryContext(ctx, query, "g", name)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var account string
		if err := rows.Scan(&account); err != nil {
			return nil, err
		}
		roles = append(roles, account)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return roles, nil
}

func(cec *CasbinEnforcerConfig) ListAllPolicies() [][]string {
	policies := cec.Enforcer.GetPolicy()
	return policies
}

func(cec *CasbinEnforcerConfig) ListRoles() []string {
	// roles := cec.Enforcer.GetAllRoles()
	roles := cec.Enforcer.GetAllSubjects()
	return roles
}

func(cec *CasbinEnforcerConfig) ListMembers() []string {
	members := make([]string, 0)
	gps := cec.Enforcer.GetGroupingPolicy()
	for _, item := range gps {
		members = append(members, item[0])
	}
	list := removeDuplicated(members)
	return list
}

func(cec *CasbinEnforcerConfig) GetMembersWithRole(name string) ([]string, error) {
	members, err := cec.Enforcer.GetUsersForRole(name)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func(cec *CasbinEnforcerConfig) GetRolesOfMember(name string) ([]string, error) {
	roles, err := cec.Enforcer.GetRolesForUser(name)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func(cec *CasbinEnforcerConfig) AddPolicy(cp *CasbinPolicy) error {
	if ok, err := cec.Enforcer.AddPolicy(cp.Sub, cp.Obj, cp.Act); !ok {
		if err != nil {
			return err
		}
		return fmt.Errorf("This policy already exists.")
	}
	cec.Enforcer.LoadPolicy()
	return nil
}

func(cec *CasbinEnforcerConfig) AddGroupingPolicy(cgp *CasbinGroupingPolicy) error {
	if ok, err := cec.Enforcer.AddGroupingPolicy(cgp.Member, cgp.Role); !ok {
		if err != nil {
			return err
		}
		return fmt.Errorf("This grouping policy already exists.")
	}
	cec.Enforcer.LoadPolicy()
	return nil
}

func(cec *CasbinEnforcerConfig) DeletePolicy(cp *CasbinPolicy) error {
	if ok, err := cec.Enforcer.RemovePolicy(cp.Sub, cp.Obj, cp.Act); !ok {
		if err != nil {
			return err
		}
		return fmt.Errorf("This policy not exists.")
	}
	cec.Enforcer.LoadPolicy()
	return nil
}

func(cec *CasbinEnforcerConfig) DeleteGroupingPolicy(cgp *CasbinGroupingPolicy) error {
	if ok, err := cec.Enforcer.RemoveGroupingPolicy(cgp.Member, cgp.Role); !ok {
		if err != nil {
			return err
		}
		return fmt.Errorf("This grouping policy not exists.")
	}
	cec.Enforcer.LoadPolicy()
	return nil
}

func(cec *CasbinEnforcerConfig) DeleteRole(name string) error {
	if ok, err := cec.Enforcer.RemoveFilteredPolicy(0, name); !ok {
		if err != nil {
			return err
		}
		return fmt.Errorf("This policy (role) not exists.")
	}
	if _, err := cec.Enforcer.RemoveFilteredGroupingPolicy(1, name); err != nil {
		return err
	}
	cec.Enforcer.LoadPolicy()
	return nil
}

func(cec *CasbinEnforcerConfig) DeleteMemeber(name string) error {
	if ok, err := cec.Enforcer.RemoveFilteredGroupingPolicy(0, name); !ok {
		if err != nil {
			return err
		}
		return fmt.Errorf("This groupiing policy (member) not exists.")
	}
	cec.Enforcer.LoadPolicy()
	return nil
}

func removeDuplicated(list []string) []string {
	allKeys := make(map[string]bool)
	result := make([]string, 0)
	for _, item := range list {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			result = append(result, item)
		}
	}
	return result
}
