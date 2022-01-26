<template>
  <div style="display:inline-block;vertical-align: top;height: 60px;background: #475285">

    <!--  <div style="display:flex;justify-content: center;align-items: center;height: 60px;">
      <select-link style="display:flex;justify-content: center;align-items: center;height: 60px" />
    </div>-->

    <el-scrollbar wrap-class="scrollbar-wrapper">

      <el-menu :default-active="activeMenu" class="el-menu-demo" mode="horizontal" background-color="#475285" text-color="#fff" active-text-color="#ffd04b">
        <sidebar-item v-for="route in permission_routes" :key="route.path" :item="route" :base-path="route.path" />
      </el-menu>

    </el-scrollbar>
  </div>
</template>

<script>

import SidebarItem from './SidebarItem'
import variables from '@/styles/variables.scss'
import SelectLink from '@/components/SelectLink'
import logo from '@/assets/index.ico'

export default {
  components: { SidebarItem, SelectLink },
  data() {
    return {
      logo: logo
    }
  },
  computed: {
    ...Vuex.mapGetters([
      'permission_routes',
      'sidebar'
    ]),
    activeMenu() {
      const route = this.$route
      if (route.meta.dynamic) {
        return route.matched[1].path
      }

      const { meta, path } = route
      // if set path, the sidebar will highlight the path you set
      if (meta.activeMenu) {
        return meta.activeMenu
      }
      return path
    },
    showLogo() {
      return this.$store.state.settings.sidebarLogo
    },
    variables() {
      return variables
    },

    isCollapse() {
      return !this.sidebar.opened
    }
  },
  methods: {
    toggleSideBar() {
      this.$store.dispatch('app/toggleSideBar')
    },
    async logout() {
      await this.$store.dispatch('user/logout')
      this.$router.push(`/login?redirect=${this.$route.fullPath}`)
    }
  }
}
</script>
<style>

</style>
