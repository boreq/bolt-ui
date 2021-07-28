<template>
    <div class="browse">
        <div class="top-bar">
            <a class="main-header" @click="onHeaderClick">Bolt UI</a>

            <ul v-if="selected">
                <li v-for="key in selected" :key="key.hex">
                    <key :k="key"></key>
                </li>
            </ul>
        </div>
        <div class="wrapper">
            <tree :path="path" :selected="selected"
                v-for="path in visiblePaths" :key="treeKey(path)"
                @entry="onEntry(path, $event)"></tree>

            <value :entry="selectedEntry" v-if="selectedEntry"></value>
        </div>
    </div>
</template>
<script lang="ts" src="./Browse.ts"></script>
<style scoped lang="scss" src="./Browse.scss"></style>
