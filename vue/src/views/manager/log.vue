<template>
  <div id="components-layout-demo-basic">
    <a-layout>
      <a-layout-sider style="padding-left: 20px;padding-top: 20px;min-width: 340px">
        <el-card style="min-width: 280px">
          <div
            style="height: 64px;font-weight: 500;font-size: 16px;border-bottom: 1px solid #e6e6e6; padding: 20px 9px 20px 9px;"
          >
            埋点管理
          </div>
          <el-menu :default-active="refreshtTab" :active="refreshtTab" @select="selectItems">
            <el-menu-item index="sbtj">
              <i class="el-icon-setting" />
              <span slot="title">上报统计</span>
            </el-menu-item>
            <el-menu-item index="sssj">
              <i class="el-icon-magic-stick" />
              <span slot="title">实时数据</span>
            </el-menu-item>
            <el-menu-item index="debugmodel">
              <i class="el-icon-user" />
              <span slot="title">Debug模式</span>
            </el-menu-item>
          </el-menu>
        </el-card>
      </a-layout-sider>
      <a-layout>
        <a-layout-content style="padding-top: 20px;padding-right: 20px">
          <track-data v-if="refreshtTab == 'sbtj'" />
          <real-time v-if="refreshtTab == 'sssj'" />
          <debug v-if="refreshtTab == 'debugmodel'" />
        </a-layout-content>
      </a-layout>
    </a-layout>
    <back-to-top />
  </div>
</template>
<script>
export default {
  name: 'Tag',
  components: {
    'Debug': () => import('@/views/manager/components/debug'),
    'RealTime': () => import('@/views/manager/components/realTime'),
    'TrackData': () => import('@/views/manager/components/TrackData'),
    BackToTop: () => import('@/components/BackToTop/index')
  },
  data() {
    return {
      tab: 'sbtj',
      tabArr: ['sbtj', 'sssj', 'debugmodel']
    }
  },
  computed: {
    refreshtTab: {
      get() {
        if (this.$store.state.baseData.RefreshTab == '' || this.tabArr.indexOf(this.$store.state.baseData.RefreshTab) == -1) {
          return this.tab
        }
        this.tab = this.$store.state.baseData.RefreshTab
        return this.$store.state.baseData.RefreshTab
      },
      set(val) {
        this.$store.dispatch('baseData/SETRefreshTab', val)
        this.tab = val
      }
    }
  },
  methods: {
    handleSizeChange() {

    },
    handleCurrentChange() {

    },
    selectItems(index) {
      this.refreshtTab = index
    }
  }
}
</script>

<style scoped src="@/styles/log.css"/>
