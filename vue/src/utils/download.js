
import FileSaver from 'file-saver'
import XLSX from 'xlsx'

export function elTable2Excel(_this, ref, fileName) {
  try {
    const $e = _this.$refs[ref].$el
    let $table = $e.querySelector('.el-table__fixed')
    if (!$table) {
      $table = $e
    }

    const wb = XLSX.utils.table_to_book($table, { raw: true })
    const wbout = XLSX.write(wb, { bookType: 'xlsx', bookSST: true, type: 'array' })
    FileSaver.saveAs(
      new Blob([wbout], { type: 'application/octet-stream' }),
      `${fileName}.xlsx`,
    )
  } catch (e) {
    if (typeof console !== 'undefined') console.error(e)
  }
}
