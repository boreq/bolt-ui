<template>
    <div class="alert" :class="kind">
        <slot></slot>
    </div>
</template>
<script lang="ts" src="./Alert.ts"></script>
<style scoped lang="scss" src="./Alert.scss"></style>
