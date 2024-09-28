<script lang="ts">
    import { page } from '$app/stores';
    import PasswordReset from './password-reset.svelte';
    const token: string | null = $page.url.searchParams.get('password_reset_token');
</script>

<html lang="en">
    <PasswordReset resetToken={token}></PasswordReset>
</html>
