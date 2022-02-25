<template>
  <div style="display:flex;justify-content:space-between">
    <div class="content_xwl">
      <div class="header_xwl" style="background: white">
        <div class="root_xwl">
          <div class="main_xwl">
            <a-tooltip placement="right" style="cursor: pointer">
              <template slot="title">
                <span>查看用户访问过的所有事件</span>
              </template>
              <span class="title_xwl" style="color: #202d3f">&nbsp;&nbsp;用户事件详情 <a-icon type="question-circle" />
              </span>
            </a-tooltip>
          </div>
          <div class="actions_xwl" />
        </div>
      </div>
      <split-pane :min-percent="0" :default-percent="22" split="vertical" @resize="onResize">
        <template slot="paneL">
          <div
            id="scollL"
            style="height: 95%;width: 100px;display: inline-block; height: 100%;vertical-align: top;width: 100%;background: white;"
          >
            <div
              style="width: 100%;height:  50px;margin-top: 0px;z-index: 10000;border-bottom: 1px solid #f0f2f5;background: white;display: flex;align-items: center;justify-content: center"
            >
              <div style="color: #303133;font-weight: bolder">用户属性</div>
            </div>
            <div style="width: 100%;height: calc(100% - 190px); overflow-x: hidden; overflow-y: auto;">
              <el-form style="margin-left: 30px;margin-top: 10px" label-position="left" inline class="user-attr-pannel">
                <el-form-item
                  v-for="(v,k,index) in this.userAttrForm"
                  :key="index"
                  :label="userAttrDescMap[k].concat('：')"
                >
                  <el-tag>{{ v }}</el-tag>
                </el-form-item>
              </el-form>
            </div>

            <div
              style="width: 100%;height:  50px;margin-bottom: 0px;z-index: 10000;border-top: 1px  solid #f0f2f5;background: white;display: flex;align-items: center;justify-content: center"
            >
              <el-button type="primary" @click="prevUser">&lt;前一个用户</el-button>
              <el-button type="primary" @click="nextUser">后一个用户&gt;</el-button>
            </div>
          </div>
        </template>
        <template slot="paneR">
          <a-spin tip="计算中..." :spinning="spinning">
            <div class="spin-content">
              <user-info-res :userId="uid" />
            </div>
          </a-spin>
        </template>
      </split-pane>
    </div>
  </div>
</template>

<script>
import { debounce } from 'lodash'
import draggable from 'vuedraggable'
import { UserList } from '@/api/analysis'

export default {
  name: 'UserInfo',
  components: {
    draggable,
    'UserInfoRes': () => import('@/views/user-analysis/components/UserInfoRes')
  },
  data() {
    return {
      userAttrForm: {},
      userAttrDescMap: {},
      spinning: false,
      debounceHandleSizeChange: undefined
    }
  },
  computed: {
    uid() {
      return this.$route.params.uid
    }
  },
  async beforeMount() {
    this.debounceHandleSizeChange = debounce(this.refreshRes, 500)
  },
  async mounted() {
    await this.init()
  },
  methods: {

    prevUser() {
      let index = Number(this.$route.params.index)

      const ui = this.$store.state.baseData.ui
      if (ui.length == 0) {
        this.$message({
          type: 'error',
          message: '用户群数量为0'
        })
        return
      }

      if (index == 0) {
        this.$message({
          type: 'error',
          message: '已是最前面一个用户'
        })
        return
      }

      index = index - 1

      this.$router.push({ path: '/user-analysis/user_info/' + ui[index] + '/' + index })
    },
    nextUser() {
      let index = Number(this.$route.params.index)

      const ui = this.$store.state.baseData.ui
      if (ui.length == 0) {
        this.$message({
          type: 'error',
          message: '用户群数量为0'
        })
        return
      }
      if (index == ui.length - 1) {
        this.$message({
          type: 'error',
          message: '已是最后一个用户'
        })
        return
      }

      index = index + 1

      this.$router.push({ path: '/user-analysis/user_info/' + ui[index] + '/' + index })
    },
    refreshRes() {
      this.$nextTick(() => {
        this.eventResShow = true
      })
    },
    onResize() {
      this.eventResShow = false
      this.debounceHandleSizeChange()
    },

    async init() {
      const uid = this.$route.params.uid
      this.propMap = {}
      const form = {}
      form['appid'] = this.$store.state.baseData.EsConnectID
      form['ui'] = [uid]
      const res = await UserList(form)
      if (res.code != 0) {
        this.$message({
          type: 'error',
          message: res.msg
        })
        return
      }
      if (res.data.alldata == null) {
        res.data.alldata = []
      }
      if (res.data.alldata.length == 0) {
        this.$message({
          type: 'error',
          message: '该用户不存在'
        })
        return
      }
      this.userAttrForm = res.data.alldata[0]
      this.userAttrDescMap = res.data.propMap
    }
  }

}
</script>

<style scoped src="@/styles/funnel.css"/>

<style>
.eventNameDisplayInput .ant-input {
  resize: none;
  border: none;
}

.eventNameDisplayInput .ant-input:focus {
  border: none;
  box-shadow: none;
}

.user-attr-pannel {
  font-size: 0;
}

.user-attr-pannel label {
  color: #99a9bf;
}

.user-attr-pannel .el-form-item {
  margin-right: 0;
  margin-bottom: 0;
  width: 100%;
}
</style>
