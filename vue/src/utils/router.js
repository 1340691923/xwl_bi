/* Layout */
import Layout from '@/layout'

// 动态路由列表
export const asyncRoutes = [

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
        path: 'event/:id',
        component: 'views/behavior-analysis/event',
        name: 'event',
        meta: {
          title: '事件分析',
          dynamic: true,
          icon: 'el-icon-data-line'
        }
      },

      {
        path: 'retention/:id',
        component: 'views/behavior-analysis/retention',
        name: 'retention',
        meta: {
          title: '留存分析',
          dynamic: true,
          icon: 'el-icon-data-analysis'
        }
      },
      {
        path: 'funnel/:id',
        component: 'views/behavior-analysis/funnel',
        name: 'funnel',
        meta: {
          title: '漏斗分析',
          dynamic: true,
          icon: 'el-icon-data-board'
        }
      },
      {
        path: 'trace/:id',
        component: 'views/behavior-analysis/trace',
        name: 'trace',
        meta: {
          title: '智能路径分析',
          dynamic: true,
          icon: 'el-icon-bicycle'
        }
      }
    ]
  },
  {
    path: '/user-analysis',
    component: 'layout',
    redirect: '/user-analysis/attr',
    alwaysShow: false,
    meta: {
      title: '用户分析',
      icon: 'el-icon-pie-chart'
    },
    children: [
      {
        path: 'attr/:id',
        component: 'views/user-analysis/index',
        name: 'attr',
        meta: {
          title: '用户属性分析',
          dynamic: true,
          icon: 'el-icon-s-custom'
        }
      },
      {
        path: 'group',
        component: 'views/user-analysis/group',
        name: 'group',
        meta: {
          title: '用户分群',
          icon: 'el-icon-user'
        }
      },
      {
        isInside: true,
        path: 'user_list',
        component: 'views/user-analysis/user_list',
        name: 'user_list',
        meta: {
          title: '用户列表',
          icon: 'el-icon-user-solid'
        }
      },
      {
        isInside: true,
        path: 'user_info/:uid/:index',
        component: 'views/user-analysis/user_info',
        name: 'user_info',
        meta: {
          title: '用户事件详情',
          dynamic: true,
          icon: 'el-icon-s-custom'
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
          icon: 'el-icon-s-management'
        }
      },
      {
        path: 'log',
        component: 'views/manager/log',
        name: 'log',
        meta: {
          title: '埋点管理',
          icon: 'el-icon-notebook-1'
        }
      }
    ]
  },
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
      },
      {
        path: 'operater_log',
        component: 'views/permission/operater_log',
        name: 'operater_log',
        meta: {
          title: '操作日志列表',
          icon: 'el-icon-s-order'
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

// 路由组件 映射 map
export const RoutesComponentmaps = {
  'layout': Layout,
  'views/dashboard/index': () => import('@/views/dashboard/index'), // 主页
  'views/permission/role': () => import('@/views/permission/role'),
  'views/permission/operater_log': () => import('@/views/permission/operater_log'),
  'views/permission/user': () => import('@/views/permission/user'),
  'views/behavior-analysis/event': () => import('@/views/behavior-analysis/event'),
  'views/behavior-analysis/retention': () => import('@/views/behavior-analysis/retention'),
  'views/behavior-analysis/funnel': () => import('@/views/behavior-analysis/funnel'),
  'views/behavior-analysis/trace': () => import('@/views/behavior-analysis/trace'),
  'views/user-analysis/index': () => import('@/views/user-analysis/index'),
  'views/user-analysis/group': () => import('@/views/user-analysis/group'),
  'views/user-analysis/user_info': () => import('@/views/user-analysis/user_info'),
  'views/user-analysis/user_list': () => import('@/views/user-analysis/user_list'),
  'views/manager/event': () => import('@/views/manager/event'),
  'views/manager/log': () => import('@/views/manager/log'),
  'views/app/index': () => import('@/views/app/index')

}

