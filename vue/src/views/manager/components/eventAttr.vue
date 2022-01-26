<template>
  <div>
    <el-dialog
      v-if="dialogVisible"
      width="80%"
      :close-on-click-modal="false"
      :visible.sync="dialogVisible"
      title="事件列表"
      @close="dialogVisible = false"
    >
      <a-input-search v-model="input2" placeholder="请输入搜索" allow-clear style="width: 200px" />
      <meta-attr style="margin-top: 10px" :attr="attr" :input="input2" />
    </el-dialog>

    <page-table
      v-if="tableShow"
      ref="pagetable"
      :connect-loading="loading"
      :table-list="tableData"
      :table-info="tableInfo"
      :input="input"
    >
      <el-table-column slot="operate" label="属性名" align="center" width="300">
        <template slot-scope="scope">

          <template v-if="typ == 1">
            {{ scope.row.attribute_name }}
          </template>
          <template v-else-if="eventName != ''">
            {{ scope.row.attribute_name }}
          </template>
          <template v-else>
            <a style="color: #6bb8ff" @click="openDialog(scope.row.attribute_name)">{{ scope.row.attribute_name }}</a>
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

      <el-table-column slot="operate" label="数据类型" align="center" prop="data_type" sortable width="200">
        <template slot-scope="scope">
          <el-tag :style="{'color':colorlists[scope.row.data_type]}">{{ scope.row.data_type_format }}</el-tag>
        </template>
      </el-table-column>

      <el-table-column slot="operate" label="属性类型" prop="attribute_type" sortable align="center" width="200">
        <template slot-scope="scope">
          <el-tag v-if="scope.row.attribute_type == 1" type="warning">预置属性</el-tag>
          <el-tag v-if="scope.row.attribute_type == 2">自定义属性</el-tag>
        </template>
      </el-table-column>
      <el-table-column slot="operate" label="是否显示" prop="attribute_type" sortable align="center" width="200">
        <template slot-scope="scope">
          <template v-if="scope.row.attribute_type != 1">
            <el-tag v-if="scope.row.status == 1">
              显示
            </el-tag>
            <el-tag v-else type="danger">
              隐藏
            </el-tag>
          </template>
        </template>
      </el-table-column>

      <el-table-column
        slot="operate"
        fixed="right"
        label="操作"
        width="200"
        align="center"
      >
        <template slot-scope="scope">
          <template v-if="scope.row.attribute_type != 1">
            <i
              v-if="scope.row.isEdit"
              class="el-icon-check"
              style="cursor: pointer;color: orangered;font-weight: bolder"
              @click="saveShowName(scope.row.index,scope.row.attribute_name,scope.row.show_name)"
            />
            <i
              v-if="!scope.row.isEdit"
              class="el-icon-edit"
              style="cursor: pointer;color: blue;font-weight: bolder"
              @click="editShowName(scope.row.index)"
            />
            <a-button
              v-if="scope.row.status == 1"
              type="link"
              style="color: red"
              icon="eye-invisible"
              @click="changeStatus(scope.row.attribute_name,0)"
            />
            <a-button
              v-if="scope.row.status == 0"
              type="link"
              icon="eye"
              @click="changeStatus(scope.row.attribute_name,1)"
            />
          </template>
        </template>
      </el-table-column>
    </page-table>
  </div>
</template>

<script>
import { AttrManager, AttrManagerByMeta, UpdateAttrInvisible, UpdateAttrShowName } from '@/api/metadata'

export default {
  name: 'EventAttr',
  components: {
    'PageTable': () => import('@/components/PageTable'),
    'MetaAttr': () => import('@/views/manager/components/metaAttr')
  },
  props: {
    input: {
      type: String,
      default: ''
    },
    typ: {
      type: Number,
      default: 2
    },
    eventName: {
      type: String,
      default: ''
    }

  },
  data() {
    return {
      colorlists: [

        'rgb(229, 0, 19)',

        'rgb(206,194,28)',

        'rgb(0,161,233)',

        'rgb(109,185,45)',

        'rgb(166,0,130)',

        'rgb(237,108,0)',

        'rgb(240, 28, 131)',

        'rgb(84, 21, 226)',

        'rgb( 0,128,0)',

        'rgb( 255,69,0)',

        'rgb( 255,165,0)',

        'rgb( 178,34,34)',

        'rgb( 255,0,255)',

        'rgb(65,105,225)',

        'blueviolet'

      ],
      loading: false,
      tableInfo: [{ slot: 'operate' }],
      tableData: [],
      tableShow: true,
      dialogVisible: false,
      input2: '',
      attr: ''
    }
  },

  mounted() {
    this.searchData()
  },

  methods: {
    async changeStatus(name, status) {
      const form = { 'appid': this.$store.state.baseData.EsConnectID, attribute_source: this.typ }
      form['attribute_name'] = name
      form['status'] = status
      const res = await UpdateAttrInvisible(form)
      if (res.code != 0) {
        this.$message({
          offset: 60,

          type: 'error',
          message: res.msg
        })

        return
      }

      this.$message({
        offset: 60,

        type: 'success',
        message: res.msg
      })
      this.searchData()
    },
    openDialog(attr) {
      this.attr = attr
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
    async saveShowName(index, attribute_name, showName) {
      const res = await UpdateAttrShowName({
        'appid': this.$store.state.baseData.EsConnectID,
        show_name: showName,
        attribute_name: attribute_name,
        typ: this.typ
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
      var res = {}
      if (this.eventName == '') {
        res = await AttrManager({ 'appid': this.$store.state.baseData.EsConnectID, typ: this.typ })
      } else {
        res = await AttrManagerByMeta({
          'appid': this.$store.state.baseData.EsConnectID,
          typ: this.typ,
          event_name: this.eventName
        })
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
