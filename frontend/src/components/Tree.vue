<template>
    <div class="tree" ref="tree" @scroll="onScroll">
        <spinner v-if="loadingPrevious" class="previous-spinner"></spinner>
        <entries :entries="tree.entries" :selected="selectedInTree" 
            @entry="onEntry($event)" v-if="tree">
        </entries>
        <spinner v-if="loadingNext" class="next-spinner"></spinner>

        <spinner v-if="!tree" class="main-spinner"></spinner>
    </div>
</template>
<script lang="ts" src="./Tree.ts"></script>
<style scoped lang="scss" src="./Tree.scss"></style>
