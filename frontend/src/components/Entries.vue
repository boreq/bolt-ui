<template>
    <div class="entries">
        <ul>
            <li v-for="entry in entries" :key="entry.key.hex">
                <a @click="onClick(entry)" :class="{ selected: selected === entry }">
                    <span class="icon">
                        <i v-if="entry.bucket" class="fas fa-folder"></i>
                        <i v-else class="fas fa-file"></i>
                    </span>

                    <key :k="entry.key"></key>
                </a>
            </li>
        </ul>
        <div class="empty-message" v-if="isEmpty">
            This bucket is empty.
        </div>
    </div>
</template>
<script lang="ts" src="./Entries.ts"></script>
<style scoped lang="scss" src="./Entries.scss"></style>
