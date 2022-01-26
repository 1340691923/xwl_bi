// 自动生成 路由映射文件 脚本
var asyncRoutes = [
  {
    path: '/permission',
    component: 'layout',
    redirect: '/permission/role',
    alwaysShow: true,
    meta: {
      title: '权限',
      icon: 'el-icon-user-solid'
    },
    children: [
      {
        path: 'role',
        component: 'views/permission/role',
        name: 'RolePermission',
        meta: {
          title: '角色管理',
          icon: 'el-icon-s-check'
        }
      },
      {
        path: 'user',
        component: 'views/permission/user',
        name: 'user',
        meta: {
          title: '用户管理',
          icon: 'el-icon-user'
        }
      }
    ]
  },
  {
    path: '/behavior-analysis',
    component: 'layout',
    redirect: '/behavior-analysis/index',
    alwaysShow: false,
    meta: {
      title: '行为分析',
      icon: 'el-icon-link'
    },
    children: [
      {
        path: 'event',
        component: 'views/behavior-analysis/event',
        name: 'event',
        meta: {
          title: '事件分析',
          icon: 'el-icon-link'
        }
      },
      {
        path: 'retention',
        component: 'views/behavior-analysis/retention',
        name: 'retention',
        meta: {
          title: '留存分析',
          icon: 'el-icon-link'
        }
      },
      {
        path: 'funnel',
        component: 'views/behavior-analysis/funnel',
        name: 'funnel',
        meta: {
          title: '漏斗分析',
          icon: 'el-icon-link'
        }
      }
    ]
  },
  {
    path: '/user-analysis',
    component: 'layout',
    redirect: '/user-analysis/index',
    alwaysShow: false,
    meta: {
      title: '用户分析',
      icon: 'el-icon-pie-chart'
    },
    children: [
      {
        path: 'index',
        component: 'views/user-analysis/index',
        name: 'index',
        meta: {
          title: '属性分析',
          icon: 'el-icon-pie-chart'
        }
      }
    ]
  },
  {
    path: '/manager',
    component: 'layout',
    redirect: '/manager/event',
    alwaysShow: false,
    meta: {
      title: '数据管理',
      icon: 'el-icon-edit'
    },
    children: [
      {
        path: 'event',
        component: 'views/manager/event',
        name: 'event',
        meta: {
          title: '事件管理',
          icon: 'el-icon-search'
        }
      },
      {
        path: 'index',
        component: 'views/manager/event_prototies',
        name: 'event_prototies',
        meta: {
          title: '事件属性管理',
          icon: 'el-icon-search'
        }
      },
      {
        path: 'index',
        component: 'views/manager/event_prototies',
        name: 'event_prototies',
        meta: {
          title: '用户属性管理',
          icon: 'el-icon-search'
        }
      }
    ]
  },
  {
    path: '/app',
    component: 'layout',
    children: [
      {
        path: 'app',
        component: 'views/app/index',
        name: 'index',
        meta: { title: '应用管理', icon: 'el-icon-s-goods' }
      }
    ]
  }
]

function filterAsyncRoutes(routes) {
  routes.forEach(route => {
    const tmp = { ...route }
    if (tmp.children) {
      tmp.children = filterAsyncRoutes(tmp.children)
    }
    if (tmp.component != 'layout') {
      console.log("'" + tmp.component + "':() => import('@/" + tmp.component + "'),")
    }
  })
}

filterAsyncRoutes(asyncRoutes)

