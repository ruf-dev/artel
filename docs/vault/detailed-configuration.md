# Detailed Configuration

Everything you might need to set up or adjust after you've created your first vault: managing the
vault itself, sharing it with others, using Fast Setup, and choosing how you sign in.

## Managing a Vault

Every vault appears as its own card in your Artel dashboard. Each card shows the vault's name and a
status indicator (for example, "active" while it's ready to use, or "provisioning" for the brief moment
it's still being set up).

### Renaming a Vault

Open the vault's card and choose the rename option, then type in the new name. This only changes the
label you see — your notes and sync connection are unaffected.

### Deleting a Vault

Open the vault's card and choose delete. This permanently removes the vault and everything stored in
it, so make sure you no longer need the notes inside before confirming. Deleting a vault does not affect
any other vaults you own.

## Sharing a Vault with Teammates

You don't have to keep a vault to yourself. You can invite other Artel users to work in the same
notebook, without ever giving out your own password.

1. Open the vault you want to share.
2. Add a teammate by their Artel username.
3. They get their own private access into that same vault — separate from yours, and revocable on its
   own.

Because each person has their own individual access, you can remove a single teammate later without
affecting anyone else, and no one ever needs to know or reuse your credentials.

This is the same approach a team lead would use to build a shared knowledge base: one vault, several
members, no shared logins.

## Fast Setup, In Detail

Fast Setup exists to skip the fiddly part of connecting Obsidian to your vault. Instead of manually
installing a plugin and typing in sync settings, you get a ready-made block of instructions.

- Every vault card has a **Fast Setup** option.
- Clicking it shows you a pre-written prompt built specifically for that vault.
- Click to copy the prompt, then paste it into Claude (or another AI assistant you're using). It reads
  the prompt and configures your Obsidian sync settings for you.
- You can use Fast Setup as many times as you like — for example, when setting up a new device.

If you'd rather not use an AI assistant for this step, you can still copy your vault's sync address
directly from the card and enter it into Obsidian's sync plugin by hand (see
[quickstart.md](quickstart.md) for the manual steps).

## Signing In

Artel supports two ways to access your account:

### Email and Password

Register with your email address and a password. Use the same combination to log back in later.

### Sign in with Telegram

If you have a Telegram account, you can sign in with it directly — no separate password to create or
remember. This is a good option if you'd rather not manage another password.

### Staying Signed In

However you sign in, your session automatically expires after a while as a security measure. If that
happens, you'll simply be asked to sign in again — none of your vaults or notes are affected.

## A Note on Administration

Larger deployments of Artel have a separate, restricted **Admin Panel** used by administrators to
manage the underlying infrastructure that vaults run on. This is not part of the regular vault
experience — as a normal user or team lead, you'll never need to visit it; creating, syncing, and
sharing vaults works the same regardless of what's happening behind the scenes.
