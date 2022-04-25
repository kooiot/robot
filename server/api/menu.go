package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kooiot/robot/server/common"
)

type MenuRouter struct {
}

var menuRouter MenuRouter

var menuList = []common.Menu{
	{Path: "/dashboard", Name: "Dashboard", Component: "LAYOUT", Redirect: "/dashboard/analysis", Meta: common.MenuMeta{Title: "routes.dashboard.dashboard", HideChildrenInMenu: true, Icon: "bx:bx-home"},
		Children: []common.Menu{
			{Path: "/analysis", Name: "Analysis", Component: "/dashboard/analysis/index", Meta: common.MenuMeta{Title: "routes.dashboard.analysis", HideMenu: true, HideBreadcrumb: true, Icon: "bx:bx-home", CurrentActiveMenu: "/dashboard"}},
			{Path: "/workbench", Name: "Workbench", Component: "/dashboard/workbench/index", Meta: common.MenuMeta{Title: "routes.dashboard.workbench", HideMenu: true, HideBreadcrumb: true, Icon: "bx:bx-home", CurrentActiveMenu: "/dashboard"}},
		},
	},
	{Path: "/permission", Name: "Permission", Component: "LAYOUT", Redirect: "/permission/front/page", Meta: common.MenuMeta{Title: "routes.demo.permission.permission", Icon: "carbon:user-role"},
		Children: []common.Menu{
			{Path: "back", Name: "PermissionBackDemo", Meta: common.MenuMeta{Title: "routes.demo.permission.back"},
				Children: []common.Menu{
					{Path: "page", Name: "BackAuthPage", Component: "/demo/permission/back/index", Meta: common.MenuMeta{Title: "routes.demo.permission.backPage"}},
					{Path: "btn", Name: "BackAuthBtn", Component: "/demo/permission/back/Btn", Meta: common.MenuMeta{Title: "routes.demo.permission.backBtn"}},
				},
			},
		},
	},
	{Path: "/level", Name: "Level", Component: "LAYOUT", Redirect: "/level/menu1/menu1-1", Meta: common.MenuMeta{Title: "routes.demo.level.level", Icon: "carbon:user-role"},
		Children: []common.Menu{
			{Path: "menu1", Name: "Menu1Demo", Meta: common.MenuMeta{Title: "Menu1"},
				Children: []common.Menu{
					{Path: "menu1-1", Name: "Menu11Demo", Meta: common.MenuMeta{Title: "Menu1-1"},
						Children: []common.Menu{
							{Path: "menu1-1-1", Name: "Menu111Demo", Component: "/demo/level/Menu111", Meta: common.MenuMeta{Title: "Menu1-1-1"}},
						},
					},
					{Path: "menu1-2", Name: "Menu12Demo", Component: "/demo/level/Menu12", Meta: common.MenuMeta{Title: "Menu1-2"}},
				},
			},
			{Path: "menu2", Name: "Menu2Demo", Component: "/demo/level/Menu2", Meta: common.MenuMeta{Title: "Menu2"}},
		},
	},
	{Path: "/system", Name: "System", Component: "LAYOUT", Redirect: "/system/account", Meta: common.MenuMeta{Title: "routes.demo.system.moduleName", Icon: "ion:settings-outline"},
		Children: []common.Menu{
			{Path: "account", Name: "AccountManagement", Component: "/demo/system/account/index", Meta: common.MenuMeta{Title: "routes.demo.system.account", IgnoreKeepAlive: true}},
			{Path: "account_detail/:id", Name: "AccountDetail", Component: "/demo/system/account/AccountDetail", Meta: common.MenuMeta{Title: "routes.demo.system.account_detail", HideMenu: true, CurrentActiveMenu: "/system/account", IgnoreKeepAlive: true}},
			{Path: "role", Name: "RoleManagement", Component: "/demo/system/role/index", Meta: common.MenuMeta{Title: "routes.demo.system.role", IgnoreKeepAlive: true}},
			{Path: "menu", Name: "MenuManagement", Component: "/demo/system/menu/index", Meta: common.MenuMeta{Title: "routes.demo.system.menu", IgnoreKeepAlive: true}},
			{Path: "dept", Name: "DeptManagement", Component: "/demo/system/dept/index", Meta: common.MenuMeta{Title: "routes.demo.system.dept", IgnoreKeepAlive: true}},
			{Path: "changePassword", Name: "ChangePassword", Component: "/demo/system/password/index", Meta: common.MenuMeta{Title: "routes.demo.system.password", IgnoreKeepAlive: true}},
		},
	},
	{Path: "/link", Name: "Link", Component: "LAYOUT", Meta: common.MenuMeta{Title: "routes.demo.iframe.frame", FrameSrc: "https://vvbin.cn/doc-next/"},
		Children: []common.Menu{
			{Path: "doc", Name: "Doc", Meta: common.MenuMeta{Title: "routes.demo.iframe.doc", IgnoreKeepAlive: true}},
			{Path: "https://vvbin.cn/doc-next/", Name: "DocExternal", Component: "LAYOUT", Meta: common.MenuMeta{Title: "routes.demo.iframe.docExternal"}},
		},
	},
}

func (b *MenuRouter) GetMenuList(c *gin.Context) {
	OkWithData(menuList, c)
}

func (s *MenuRouter) InitRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	router := Router.Group("api")
	{
		router.GET("getMenuList", menuRouter.GetMenuList)
	}
	return router
}
