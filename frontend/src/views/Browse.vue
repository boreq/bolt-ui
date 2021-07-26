<template>
    <div class="browse">
        <div class="top-bar">
            <a class="main-header" @click="onHeaderClick">Bolt UI</a>

            <ul>
                <li v-for="key in selectedPath" :key="key.hex">
                    <key :k="key"></key>
                </li>
            </ul>
        </div>
        <div class="wrapper">
            <entries :entries="tree.entries" :selected="selectedInTree(tree)"
                v-for="tree in visibleTrees" :key="treeKey(tree)"
                @entry="onEntry(tree, $event)"></entries>

            <value :entry="selected" v-if="selected"></value>
        </div>
    </div>
</template>
<script lang="ts" src="./Browse.ts"></script>
<style scoped lang="scss" src="./Browse.scss"></style>
