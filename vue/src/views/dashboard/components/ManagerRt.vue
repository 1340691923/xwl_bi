<template>
  <div>
    <a-drawer
      title="管理看板内容"
      :width="720"
      :visible="visible"
      :body-style="{ paddingBottom: '80px' }"
      @close="onClose"
    >
      <div style="flex-direction: row;display: flex;justify-content: space-between">
        <div id="components-layout-basic" style="width: 49%;height: 500px">
          <a-layout>
            <a-layout-header style="display: flex;flex-direction: row;justify-content: space-between">
              <div>
                <a-icon type="drag" />
                调整布局
              </div>
              <div>
                报表数({{ selectRts.length }}/{{ tableData.length }})
              </div>
            </a-layout-header>
            <a-layout-content>
              <draggable v-model="selectRts" chosen-class="chosen" animation="1000">
                <a-dropdown v-for="(v,k) in selectRts" :key="Date.now()" v-drag :trigger="['contextmenu']">
                  <div
                    style="width: 49%;height: 100px;line-height: 50px;text-align: center;padding: 20px;float: left;border: 2px solid grey;"
                  >
                    {{ rtMap[v] }}
                  </div>
                  <a-menu slot="overlay">
                    <a-menu-item key="1" @click="removeRt(v,k)">
                      删除报表
                    </a-menu-item>
                  </a-menu>
                </a-dropdown>
              </draggable>
            </a-layout-content>
          </a-layout>
        </div>
        <div style="width: 49%;height: 600px">
          <div
            style="display: flex;width: 100%;padding: 5px;height: 100%;flex-direction: column;justify-content: flex-start;align-items: flex-start"
          >
            <el-select v-model="rtType" style="margin: 0px auto;width: 90%" placeholder="全部报表" @change="fliter">
              <el-option
                v-for="(v,k,index) in tableTypeMap"
                :key="index"
                :label="v"
                :value="Number(k)"
              />
            </el-select>
            <el-input
              v-model="input"
              placeholder="请输入报表名称"
              prefix-icon="el-icon-search"
              style="margin: 0px auto;margin-top: 20px;width: 90%"
              @input="fliter"
            />
            <div
              v-if="showRtList"
              v-infinite-scroll="handleInfiniteOnLoad"
              style="margin: 0px auto;margin-top: 20px;width: 90%"
              class="infinite-container"
              :infinite-scroll-disabled="busy"

              :infinite-scroll-distance="1000"
            >
              <a-list :data-source="data" size="small" :bordered="true" :split="true">
                <a-list-item v-if="!item.isAdd" slot="renderItem" slot-scope="item, index">
                  <a-tooltip placement="left" style="cursor: pointer">
                    <template slot="title">
                      <span>{{ item.remark }}</span>
                    </template>
                    <a-list-item-meta :description="tableTypeMap[item.rt_type]">
                      <a slot="title">{{ item.name }}</a>
                    </a-list-item-meta>
                  </a-tooltip>
                  <div>
                    <el-button type="warning" icon="el-icon-plus" circle @click="add2Pannel(item.id)" />
                  </div>
                </a-list-item>
                <div v-if="loading && !busy" class="demo-loading-container">
                  <a-spin />
                </div>
              </a-list>
            </div>
          </div>
        </div>
      </div>
      <div
        :style="{
          position: 'absolute',
          right: 0,
          bottom: 0,
          width: '100%',
          borderTop: '1px solid #e9e9e9',
          padding: '10px 16px',
          background: '#fff',
          textAlign: 'right',
          zIndex: 1,
        }"
      >
        <a-button :style="{ marginRight: '8px' }" @click="onClose">
          取消
        </a-button>
        <a-button type="primary" @click="go">
          应用
        </a-button>
      </div>
    </a-drawer>
  </div>
</template>

<script>
import draggable from 'vuedraggable'
import infiniteScroll from 'vue-infinite-scroll'
import { ReportTableList } from '@/api/pannel'
import { filterData } from '@/utils/table'

export default {
  directives: { infiniteScroll },
  components: {
    draggable
  },
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    value: {
      type: Array,
      default: []
    }
  },
  data() {
    return {
      showRtList: false,
      selectRts: this.value,
      input: '',
      tableData: [],
      rtType: 0,
      data: [],
      loading: false,
      busy: false,
      tableTypeMap: {
        0: '全部报表',
        1: '事件分析',
        2: '留存分析',
        3: '漏斗分析',
        4: '智能路径分析',
        5: '用户属性分析'
      },
      rtMap: {}
    }
  },
  beforeMount() {
    this.init()
  },
  methods: {
    go() {
      const selectRts = []
      for (const v of this.selectRts) {
        selectRts.push(v.toString())
      }
      this.$emit('input', selectRts)
      this.$emit('saveRt', selectRts)
      this.onClose()
    },
    async init() {
      await this.getReportTableList()
    },
    onClose() {
      this.$emit('close')
    },
    removeRt(id, index) {
      this.selectRts.splice(index, 1)

      for (const k in this.tableData) {
        const v = this.tableData[k]

        if (v.id.toString() == id) {
          this.tableData[k].isAdd = false
        }
      }
      this.fliter()
    },
    add2Pannel(id) {
      this.selectRts.push(id)
      for (const k in this.tableData) {
        const v = this.tableData[k]

        if (v.id.toString() == id) {
          this.tableData[k].isAdd = true
        }
      }
      this.fliter()
    },
    fliter() {
      this.data = filterData(this.tableData, this.input)
      this.showRtList = false
      this.$nextTick(() => {
        this.showRtList = true
      })
    },

    async getReportTableList() {
      const res = await ReportTableList({ appid: this.$store.state.baseData.EsConnectID, rt_type: this.rtType })
      if (res.code != 0) {
        this.$message({
          offset: 60,
          type: 'error',
          message: res.msg
        })
        return
      }
      if (res.data == null) {
        res.data = []
      }
      this.tableData = res.data
      this.rtMap = {}

      for (const k in this.tableData) {
        const v = this.tableData[k]
        this.rtMap[v.id] = v.name

        if (this.selectRts.indexOf(v.id.toString()) != -1) {
          this.tableData[k].isAdd = true
        } else {
          this.tableData[k].isAdd = false
        }
      }

      this.fliter()
    },
    handleInfiniteOnLoad() {

    }
  }
}
</script>
<style>
#components-layout-basic {
  text-align: center;
}

#components-layout-basic .ant-layout-header {
  color: #42546d;
  background: white;
}

.infinite-container {
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  overflow: auto;
  padding: 8px 24px;
  height: 700px;
}

.loading-container {
  position: absolute;
  bottom: 40px;
  width: 100%;
  text-align: center;
}

.actions_xwl_btn:hover {
  color: orangered;
}

</style>
