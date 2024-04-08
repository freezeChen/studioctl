<script setup lang="ts">
import {getTableColumns, getTables, TableColumnRes, TableRes} from '@/api/module/auto';
import {ref} from 'vue';


const searchOpt = ["=", "like", "between"]

const tablesRef = ref<TableRes[]>()
const selectTableRef = ref('')
const columnsRef = ref<TableColumnRes[]>()


async function initData() {
  tablesRef.value = await getTables()
  console.log(tablesRef.value);

}

initData()

async function loadTableColumns() {
  const tableColumns = await getTableColumns(selectTableRef.value)
  columnsRef.value = tableColumns


}

</script>

<template>
  <el-form>
    <el-form-item>
      <el-select v-model="selectTableRef">
        <el-option :key="v.table_name" :label="v.table_comment + v.table_name" :value="v.table_name"
                   v-for="v in tablesRef"/>
      </el-select>
    </el-form-item>
  </el-form>
  <el-button @click="loadTableColumns">确定</el-button>

  <el-table :data="columnsRef" :border="true">
    <el-table-column prop="columnName" label="字段名称"/>
    <el-table-column prop="dataType" label="数据类型"/>
    <el-table-column prop="columnComment" label="注释"/>
    <el-table-column prop="show" label="是否显示">
      <template #default="scope">
        <input type="checkbox"/>
      </template>
    </el-table-column>
    <el-table-column prop="search" label="搜索条件">
      <template #default="scope">
        <el-select v-model="scope.row.search">
          <el-option v-for="s in searchOpt" :label="s" :value="s"/>
        </el-select>
      </template>
    </el-table-column>

  </el-table>
</template>

<style scoped></style>
