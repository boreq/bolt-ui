import { Vue } from 'vue-property-decorator';
import { Location } from 'vue-router';

export class NavigationService {

    // todo remove
    constructor(private vue: Vue) {
    }

    // todo remove
    escapeHome(): void {
        this.vue.$router.replace(
            {
                name: 'browse',
            },
        );
    }

    getBrowse(): Location {
        return {
            name: 'browse',
        };
    }

    getSettings(): Location {
        return {
            name: 'settings',
        };
    }

    getSettingsProfile(): Location {
        return {
            name: 'settings-profile',
        };
    }

    getSettingsPrivacyZones(): Location {
        return {
            name: 'settings-privacy-zones',
        };
    }

    getSettingsImport(): Location {
        return {
            name: 'settings-import',
        };
    }

    getSettingsInstance(): Location {
        return {
            name: 'settings-instance',
        };
    }

    getProfile(username: string): Location {
        return {
            name: 'profile',
            params: {
                username: username,
            },
        };
    }

    getProfileWithBefore(username: string, before: string): Location {
        return {
            name: 'profile',
            params: {
                username: username,
            },
            query: {
                before: before,
            },
        };
    }

    getProfileWithAfter(username: string, after: string): Location {
        return {
            name: 'profile',
            params: {
                username: username,
            },
            query: {
                after: after,
            },
        };
    }

    getActivity(activityUUID: string): Location {
        return {
            name: 'activity',
            params: {
                activityUUID: activityUUID,
            },
        };
    }

    getActivitySettings(activityUUID: string): Location {
        return {
            name: 'activity-settings',
            params: {
                activityUUID: activityUUID,
            },
        };
    }

    getNewPrivacyZone(): Location {
        return {
            name: 'new-privacy-zone',
        };
    }

    getPrivacyZoneSettings(privacyZoneUUID: string): Location {
        return {
            name: 'privacy-zone-settings',
            params: {
                privacyZoneUUID: privacyZoneUUID,
            },
        };
    }

    getNewActivity(): Location {
        return {
            name: 'new-activity',
        };
    }

}
