<script setup lang="ts">
import {getTableColumns, getTables, Preview, previewCode, TableRes} from '@/api/module/auto';
import {ref} from 'vue';
import 'highlight.js/styles/stackoverflow-light.css'
import 'highlight.js/lib/common';
import hljsVuePlugin from "@highlightjs/vue-plugin";


const searchOpt = ["=", "like", "between"]

const tablesRef = ref<TableRes[]>()
const selectTableRef = ref('')
const columnsRef = ref<Preview>()
const showCodeRef = ref(false)
const previewCodeRef = ref('')


async function initData() {
  tablesRef.value = await getTables()
  console.log(tablesRef.value);

}

initData()

async function loadTableColumns() {
  const tableColumns = await getTableColumns(selectTableRef.value)

  columnsRef.value = tableColumns
}

async function preview() {
  const code = await previewCode(columnsRef.value!!)

  console.log(code)
  showCodeRef.value = !showCodeRef.value;
  previewCodeRef.value = code

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
  <el-form label-width="auto">
    <el-form-item label-width="名称">

    </el-form-item>

  </el-form>

  <el-table :data="columnsRef?.fields" :border="true">
    <el-table-column prop="field_name" label="字段名称" width="140">
      <template #default="scope">
        <el-input type="text" v-model="scope.row.field_name"/>
      </template>

    </el-table-column>
    <el-table-column prop="field_zh_name" label="中文名称" width="120">
      <template #default="scope">
        <el-input type="text" v-model="scope.row.field_zh_name"/>
      </template>
    </el-table-column>
    <el-table-column prop="field_type" label="数据类型" width="180">
      <template #default="scope">
        <el-input type="text" v-model="scope.row.field_type"/>
      </template>

    </el-table-column>
    <el-table-column prop="field_comment" label="注释" width="130">
      <template #default="scope">
        <el-input type="text" v-model="scope.row.field_comment"/>
      </template>
    </el-table-column>
    <el-table-column prop="field_json" label="数据库字段" min-width="120" >
      <template #default="scope">
        <el-input type="text" v-model="scope.row.field_json"/>
      </template>
    </el-table-column>
    <el-table-column prop="show" label="是否显示">
      <template #default="scope">
        <input type="checkbox" v-model="scope.row.show"/>
      </template>
    </el-table-column>
    <el-table-column prop="require" label="编辑是否必填">
      <template #default="scope">
        <input type="checkbox" v-model="scope.row.require"/>
      </template>
    </el-table-column>

    <el-table-column prop="search_type" label="搜索条件" min-width="100">
      <template #default="scope">
        <el-select v-model="scope.row.search_type" clearable>
          <el-option v-for="s in searchOpt" :label="s" :value="s"/>
        </el-select>
      </template>
    </el-table-column>

  </el-table>

  <el-dialog v-model="showCodeRef">
    <hljsVuePlugin.component :code="previewCodeRef"/>

  </el-dialog>
  <el-button @click="preview">预览</el-button>
</template>

<style scoped></style>
