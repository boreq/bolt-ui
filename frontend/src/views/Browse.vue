<template>
    <div class="browse">
        <div class="top-bar">
            <a class="main-header" @click="onHeaderClick">Bolt UI</a>

            <ul v-if="selectedPath">
                <li v-for="key in selectedPath" :key="key.hex">
                    <key :k="key"></key>
                </li>
            </ul>
        </div>
        <div class="wrapper">
            <tree :path="path" :selected="selectedPath"
                v-for="path in visiblePaths" :key="treeKey(path)"
                @entry="onEntry(path, $event)"></tree>

            <value :entry="selectedValue" v-if="selectedValue"></value>
        </div>
    </div>
</template>
<script lang="ts" src="./Browse.ts"></script>
<style scoped lang="scss" src="./Browse.scss"></style>
