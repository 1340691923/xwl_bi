<template>
  <div id="components-layout-demo-basic">
    <a-layout>
      <a-layout-sider style="padding-left: 20px;padding-top: 20px;min-width: 340px">
        <el-card style="min-width: 280px">
          <div
            style="height: 64px;font-weight: 500;font-size: 16px;border-bottom: 1px solid #e6e6e6; padding: 20px 9px 20px 9px;"
          >
            元数据管理
          </div>
          <el-menu :default-active="refreshtTab" :active="refreshtTab" @select="selectItems">
            <el-menu-item v-for="(v,k,index) in tabObj" :key="index" :index="k">
              <i :class="v.icon" />
              <span slot="title">{{ v.title }}</span>
            </el-menu-item>
          </el-menu>
        </el-card>
      </a-layout-sider>
      <a-layout>
        <a-layout-content style="padding-top: 20px;padding-right: 20px">
          <el-card class="box-card">
            <div
              style="height: 50px;line-height: 50px;display: flex;align-items: center;justify-content: space-between;border-bottom: 1px solid #f0f2f5"
            >
              <a-tooltip placement="right" style="cursor: pointer">
                <template slot="title">
                  <span>{{ tabObj[refreshtTab].desc }}</span>
                </template>
                <span class="title_xwl" style="color: #202d3f">{{ tabObj[refreshtTab].title }}&nbsp;<a-icon
                  type="question-circle"
                /></span>
              </a-tooltip>
              <div>
                <a-input-search v-model="input" placeholder="请输入搜索" allow-clear style="width: 200px" />
              </div>
            </div>
            <meta-attr v-if="refreshtTab == 'metaAttr'" :input="input" />
            <event-attr v-if="refreshtTab == 'eventAttr'" key="1" :typ="Number(2)" :input="input" />
            <event-attr v-if="refreshtTab == 'userAttr'" key="2" :typ="Number(1)" :input="input" />
          </el-card>
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
    'EventAttr': () => import('@/views/manager/components/eventAttr'),
    'MetaAttr': () => import('@/views/manager/components/metaAttr'),
    BackToTop: () => import('@/components/BackToTop/index')
  },
  data() {
    return {
      input: '',
      tab: 'metaAttr',
      tabArr: ['metaAttr', 'eventAttr', 'userAttr'],
      tabObj: {
        'metaAttr': {
          title: '元事件',
          desc: '在该页面进行事件及事件属性的管理。包括设置事件的显示名、显示状态等；设置事件属性的显示名、显示状态、计数单位，上传维度表等功能',
          icon: 'el-icon-setting'
        },
        'eventAttr': {
          title: '事件属性',
          desc: '在该页面进行事件属性的管理。包括设置事件属性的显示名、显示状态、计数单位，上传维度表以及设置虚拟事件属性等功能',
          icon: 'el-icon-magic-stick'
        },
        'userAttr': {
          title: '用户属性',
          desc: '在该页面进行用户属性的管理。包括设置用户属性的显示名、显示状态、计数单位，上传维度表以及设置虚拟用户属性等功能',
          icon: 'el-icon-user'
        }
      }
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
    selectItems(index) {
      this.refreshtTab = index
    }
  }
}
</script>
<style scoped src="@/styles/manager-event.css"/>
