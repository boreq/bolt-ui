<template>
    <div class="value">
        <div class="header">
            <i class="fas fa-file"></i>

            <key :k="entry.key"></key>

            <div class="format-note">
                <span v-tooltip="formatTooltip">({{ format }})</span>
            </div>
        </div>

        <div class="value-string" v-if="entry.value">
            <div v-if="valuePretty">
                <div class="value-header">
                    Pretty printed
                </div>
                <pre class="value-string"><code>{{ valuePretty }}</code></pre>
            </div>

            <div class="value-header">
                Raw value as hex
            </div>
            <pre class="value-string"><code>{{ valueHex }}</code></pre>
        </div>

        <div class="value-empty" v-else>
            This value is not set.
        </div>
    </div>
</template>
<script lang="ts" src="./Value.ts"></script>
<style scoped lang="scss" src="./Value.scss"></style>
