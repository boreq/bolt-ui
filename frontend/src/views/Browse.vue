<template>
    <div class="browse">
        <div class="top-bar">
            <a class="main-header" @click="onHeaderClick">Bolt UI</a>

            <ul v-if="selectedPath && !editingSelectedPath" @click.stop="onSelectedPathClick">
                <li v-for="key in selectedPath" :key="key.hex">
                    <key :k="key"></key>
                </li>
            </ul>

            <div class="edit-path" v-if="editingSelectedPath">
                <input v-model="editedPath" class="path-input"
                    @keyup.enter="onEditPathSubmit" @click.stop>
            </div>
        </div>
        <div class="wrapper">
            <tree :path="path" :selected="selectedPath"
                v-for="(path, index) in paths" :key="treeKey(path)"
                @entry="onEntry(path, $event)" @path="onPath"
                v-show="isTreeVisible(index)"></tree>

            <value :entry="selectedValue" v-if="selectedValue"></value>
        </div>
    </div>
</template>
<script lang="ts" src="./Browse.ts"></script>
<style scoped lang="scss" src="./Browse.scss"></style>
