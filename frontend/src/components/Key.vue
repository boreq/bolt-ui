<template>
    <div class="key">
        <span v-if="k.str">
            <span class="decoration">"</span>{{ k.str }}<span class="decoration">"</span>
        </span>
        <span v-else>
            <span class="decoration">0x</span>{{ k.hex }}
        </span>
    </div>
</template>
<script lang="ts" src="./Key.ts"></script>
<style scoped lang="scss" src="./Key.scss"></style>
