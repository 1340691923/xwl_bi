<template>
  <div class="app-container">
    <el-card class="box-card" style="min-height:  850px;">
      <div
        style="height: 50px;line-height: 50px;display: flex;align-items: center;justify-content: space-between;border-bottom: 1px solid #f0f2f5"
      >

        <a-tooltip placement="right" style="cursor: pointer">
          <template slot="title">
            <span>用户列表</span>
          </template>
          <span class="title_xwl" style="color: #202d3f">&nbsp;&nbsp;用户列表 <a-icon
            type="question-circle"
          />
          </span>
        </a-tooltip>
        <div>
          <a-input-search v-model="input" placeholder="请输入搜索" style="width: 200px" />
        </div>
      </div>

      <page-table
        v-if="tableShow"
        :connect-loading="loading"
        :border="false"
        style="padding: 20px"
        :input="input"
        :limit="Number(10)"
        :table-list="tableData"
        :table-info="tableInfo"
      >
        <el-table-column slot="operate" fixed label="访客ID" align="center" sortable>
          <template slot-scope="scope">
            <a style="color: #6bb8ff" @click="lookUserInfo(scope.row['xwl_distinct_id'],scope.$index)">{{ scope.row['xwl_distinct_id'] }}</a>
          </template>
        </el-table-column>
        <el-table-column
          v-for="(v,k,index) in propMap"
          slot="operate"
          :key="index"
          :label="v"
          :prop="k"
          align="center"
          sortable
        >
          <template slot-scope="scope">
            {{ scope.row[k] }}
          </template>
        </el-table-column>
      </page-table>

    </el-card>
  </div>
</template>

<script>
import { UserList } from '@/api/analysis'

export default {
  name: 'UserList',
  components: {
    'PageTable': () => import('@/components/PageTable')
  },
  data() {
    return {
      loading: false,
      tableInfo: [{ slot: 'operate' }],
      tableData: [],
      propMap: {},
      input: '',
      tableShow: false
    }
  },
  async mounted() {
    await this.init()
  },
  methods: {
    refreshTable() {
      this.tableShow = false
      this.$nextTick(() => {
        this.tableShow = true
      })
    },
    async init() {
      this.propMap = {}
      const form = {}
      form['appid'] = this.$store.state.baseData.EsConnectID
      form['ui'] = this.$store.state.baseData.ui
      this.loading = true
      const res = await UserList(form)
      this.loading = false
      if (res.code != 0) {
        this.$message({
          offset: 60,
          type: 'error',
          message: res.msg
        })
        this.refreshTable()
        return
      }
      if (res.data.alldata == null) {
        res.data.alldata = []
      }
      this.tableData = res.data.alldata

      Object.keys(res.data.propMap).sort().forEach((key) => {
        this.propMap[key] = res.data.propMap[key]
      })

      delete this.propMap['xwl_distinct_id']
      this.refreshTable()
    },
    lookUserInfo(uid, index) {
      this.$router.push({ path: '/user-analysis/user_info/' + uid + '/' + index })
    }
  }
}
</script>

<style scoped>

</style>
