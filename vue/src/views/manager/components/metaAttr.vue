<template>
  <div>
    <el-dialog
      v-if="dialogVisible"
      width="80%"
      :close-on-click-modal="false"
      :visible.sync="dialogVisible"
      title="属性列表"
      @close="dialogVisible = false"
    >
      <a-input-search v-model="input2" placeholder="请输入搜索" allow-clear style="width: 200px" />
      <event-attr style="margin-top: 10px" :typ="Number(2)" :event-name="eventName" :input="input2" />
    </el-dialog>

    <page-table
      v-if="tableShow"
      ref="pagetable"
      :connect-loading="loading"
      :table-list="tableData"
      :table-info="tableInfo"
      :input="input"
    >
      <el-table-column slot="operate" label="事件名" align="center" width="300">
        <template slot-scope="scope">
          <template v-if="attr == ''">
            <a style="color: #6bb8ff" @click="openDialog(scope.row.event_name)">{{ scope.row.event_name }}</a>
          </template>
          <template v-else>
            {{ scope.row.event_name }}
          </template>
        </template>
      </el-table-column>
      <el-table-column slot="operate" label="显示名" align="center" width="300">
        <template slot-scope="scope">
          <div v-if="!scope.row.isEdit">
            {{ scope.row.show_name }}
          </div>
          <div v-else>
            <el-input v-model="scope.row.show_name" size="mini" />
          </div>
        </template>
      </el-table-column>
      <el-table-column slot="operate" label="昨日事件量" align="center" prop="yesterday_count" sortable />
      <el-table-column
        slot="operate"
        fixed="right"
        label="操作"
        width="200"
        align="center"
      >
        <template slot-scope="scope">
          <i
            v-if="scope.row.isEdit"
            class="el-icon-check"
            style="cursor: pointer;color: orangered;font-weight: bolder"
            @click="saveShowName(scope.row.index,scope.row.event_name,scope.row.show_name)"
          />
          <i
            v-else
            class="el-icon-edit"
            style="cursor: pointer;color: blue;font-weight: bolder"
            @click="editShowName(scope.row.index)"
          />
        </template>
      </el-table-column>
    </page-table>
  </div>
</template>

<script>
import { MetaEventList, MetaEventListByAttr, UpdateShowName } from '@/api/metadata'

export default {
  name: 'MetaAttr',
  components: {
    'PageTable': () => import('@/components/PageTable'),
    'EventAttr': () => import('@/views/manager/components/eventAttr')
  },
  props: {
    input: {
      type: String,
      default: ''
    },
    attr: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      eventName: '',
      loading: false,
      tableInfo: [{ slot: 'operate' }],
      tableData: [],
      tableShow: true,
      dialogVisible: false,
      input2: ''
    }
  },

  mounted() {
    this.searchData()
  },

  methods: {
    openDialog(eventName) {
      this.eventName = eventName
      this.dialogVisible = true
    },
    editShowName(index) {
      for (const i in this.tableData) {
        if (this.tableData[i].index == index) {
          this.tableData[i].isEdit = true
        }
      }
      this.refreshTable()
    },
    async saveShowName(index, eventName, showName) {
      const res = await UpdateShowName({
        'appid': this.$store.state.baseData.EsConnectID,
        show_name: showName,
        event_name: eventName
      })

      if (res.code != 0) {
        this.$message({
          offset: 60,

          type: 'error',
          message: res.msg
        })
        return
      }

      for (const i in this.tableData) {
        if (this.tableData[i].index == index) {
          this.tableData[i].isEdit = false
        }
      }
      this.refreshTable()
    },
    async searchData() {
      this.loading = true
      let res = {}
      if (this.attr == '') {
        res = await MetaEventList({ 'appid': this.$store.state.baseData.EsConnectID })
      } else {
        res = await MetaEventListByAttr({ 'appid': this.$store.state.baseData.EsConnectID, 'attr': this.attr })
      }

      this.loading = false
      if (res.code != 0) {
        this.$message({
          offset: 60,

          type: 'error',
          message: res.msg
        })

        return
      }

      this.tableData = []
      for (const k in res.data.list) {
        const v = res.data.list[k]
        v['index'] = k
        v['isEdit'] = false
        this.tableData.push(v)
      }

      this.refreshTable()
    },
    refreshTable() {
      this.tableShow = false
      this.$nextTick(() => {
        this.tableShow = true
      })
    }
  }
}
</script>

<style scoped>

</style>
