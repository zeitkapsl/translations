import { t } from "$lib/i18n/i18n.js";

export default {
  "en-US": {
    "en-US": "EN-US",
    "de-DE": "DE-DE",
    "de-AT": "DE-AT",
    "zeitkapsl.rights": `&copy; ${new Date().getFullYear()} zeitkapsl. All rights reserved.`,
    login: "Sign in",
    logout: "Sign out",
    register: "Register",
    email: "E-Mail address",
    password: "Password",
    "password.forgot": "Forgot password?",
    save: "Save",
    saving: "Saving",
    next: "Next",
    repeat: "repeat",
    "not a member": "Not a member?",
    settings: "Settings",
    language: "Language",
    selected: "selected",
    unselect: "unselect",
    "register.success": "Please check your mail",
    "step.upper.case": "STEP",
    "of.upper.case": "OF",
    of: "of",
    first_name: "First Name",
    last_name: "Last Name",
    name: "Name",
    "back.to.home": "Back to home",
    back: "Back",
    "error.short": "Error",
    error: "Unexpected problem occurred",
    "try.later": "Please try again later!",
    "currently.not.available": "Currently not available",
    edit: "Edit",
    updated: "Updated",
    conflict: "The item got updated",
    "conflict.detail": "We reloaded the media item!",
    elements: "Elements",
    element: "Element",
    open: "Open",
    create: "Create",
    creating: "Creating",
    cancel: "Cancel",
    cancelling: "Cancelling",
    tasks: "Tasks",
    "tasks.in.progress": "Tasks are in progress",
    processing: "Processing",
    "processing.done": "Processing done",
    "processing.done.with.errors": "Processing done with errors",
    search: "Search",
    menu: "Menu",
    info: "Info",
    remove: "Remove",
    select: "Select",
    "media.not.found": "Media not found",
    "drag.and.drop.or.click": "Drag and drop or click",
    "or.get.the.app": "or get the app",
    change: "Change",
    "thank.you": "Thank you",
    "lets.go": "Let's go!",
    retry: "Retry",
    account: "Account",
    accounts: "Accounts",
    "email.already.exists": "Email address already exists",
    "personal.information": "Personal information",
    members: "Members",
    "not.editable": "Read-only",
    add: "Add",
    "click.here": "Click here",
    "saving.changes": "Saving changes",
    "changes.saved": "Saved changes",
    syncing: "Sync in progress",
    plan: "Subscription",
    photos: "Photos",
    minimize: "Minimize Window",
    maximize: "Maximize Window",

    "collection.empty": "Your album is empty",
    "collection.empty.add": "Add something to your album",
    "collection.add": "Add to album",
    "collection.add.to": (collectionName) => `Add to album '${collectionName}'`,
    "collection.added": "Added to album",
    "collection.name": "Name of album",
    "collection.name.enter": "Enter name of album",
    "collection.delete": "Delete album",
    "collection.delete.detail":
      "Deleting the album cannot be undone. Photos and videos from the deleted album remain saved in the zeitkapsl.",
    "collection.remove": "Remove from album",
    "collection.remove.detail":
      "Photos and videos remain saved in your photo gallery.",
    "collection.goto": "Go to album",
    "collection.as.cover.media": "Select as album cover",
    "no.collections": "No Alben found",
    "no.collections.add": "Do you want to create one?",
    "no.shared.collections": "No shared album found",

    "nothing.found": "Nothing found",

    "read.only.account": "Demo account",
    "read.only.account.detail": "Therefore, this feature is not available.",

    "storage.limit.reached": "Speicherplatz voll",
    "storage.limit.reached.detail":
      "Dein Speicherplatz ist aufgebraucht. Bitte erweitere deinen Speicher.",

    "share.via.link": "Share via link",
    "share.link.create": "Create Link",
    "accept.share": "Add this shared album in my albums?",
    "copy.shared.collection": "Copy album",
    "copy.shared.collection.success": "Album copied",
    "exit.shared.collection": "Remove",
    "exit.shared.collection.deleted": "Removed",
    "share.link.valid": (validUntil) =>
      `Sharing link, valid until ${validUntil}`,
    "share.link.edit": "Sharing link",
    "share.link.allow.uploads": "Allow uploads",
    "valid.until": "Valid until",
    "share.delete.titel": "Delete sharing link",
    "share.delete.detail":
      "Deleting the link cannot be undone. Do you want to delete the sharing link?",

    waitlist: "Join our waitlist",
    "waitlist.motivation":
      "Belong to the first who access zeitkapsl. <br\> <br\> If you don't know anybody with an invite, you can sign up for the waitlist and we'll send one soon.",
    "waitlist.register": "Join now",
    "waitlist.success.title": "Thank you very much!",
    "waitlist.success.detail": "We have added you to the waiting list.",

    "register.welcome": "Choose your perfect subscription",
    "register.select.subscription.motivation":
      "Start now and discover what's possible!",
    "register.first.step.motivation": "Just a few steps and you're done!",
    "register.create.password": "Create your account",
    "register.invite.code": "Enter invitation code",
    "register.invite.code.motivation":
      "To participate you need an invitation code!",
    "register.invite.code.question": "Don't have an invitation code? Join the ",
    "register.invite.code.question.resolve": "waitlist.",
    "register.invite.code.question.resolve.final": " ",
    "register.user.password.strength": "Password strength: ",
    "register.user.password.strength.too.weak": "too weak",
    "register.user.password.strength.weak": "weak",
    "register.user.password.strength.medium": "medium",
    "register.user.password.strength.strong": "strong",
    "register.payment.information": "Choose your preferred payment method",
    "register.payment.information.details1":
      "Your payments are safe with us. You can easily quit online.",
    "register.payment.information.details2":
      "Your details are encrypted and you can change the payment method at any time.",
    "register.recovery": "Save your Secret Key",
    "register.recovery.details":
      "Your Secret Key belongs only to you. Your data is encrypted with this key. We <span class='font-bold'>don't</span> store it! <br\> <br\> If you forget your password, the only way you can recover your data is with this key. Download and store it somewhere safe.",
    "register.recovery.details.short":
      "Your Secret Key belongs only to you. Your data is encrypted with this key. We <span class='font-bold'>don't</span> store it! <br\> <br\> If you forget your password, the only way you can recover your data is with this key.",
    "register.recovery.wordlist": "Your personal recovery words",
    "register.recovery.qrcode": "QR Code",
    "register.recovery.qrcode.detail": "Scan when resetting your password.",
    "register.recovery.secret.key": "Secret key",

    "register.recovery.pdf.title": "Emergency Kit",
    "register.recovery.pdf.description":
      "If you forget your password, the only way you can recover your data and reset your password is with this secret key.\n\n**We cannot access nor recover your data for you. ** Store it somewhere safe, such as your birth certificate.",
    "register.recovery.pdf.account": "Created for",
    "register.recovery.needhelp": "Need help?",
    "register.recovery.needhelp.detail": "support@zeitkapsl.eu",
    "register.recovery.pdf.download": "Download Secret Key",
    "register.recovery.pdf.download.again": "Download again",
    "register.in_progress": "We are creating your account",
    "register.for_free": "We are happy to offer this for free :)",
    "register.now": "Register now",
    "register.payment.success.detail":
      "Thank you for your payment! Your account is ready to be used.",

    "payment.failed": "Payment failed",

    "recovery.key": "Key",
    "recovery.scan.qr.code": "Scan QR-Code",
    "recovery.scan.qr.code.stop": "Close Scanner",

    "reset.password": "Reset your password",
    "reset.password.detail":
      "Enter your email and a new password to reset the password on your account.",
    "reset.password.new": "Enter new password",
    "reset.password.do": "Reset password",

    "zeitkapsl.accept.terms.and.conditions":
      'I have read and agree with the <a href="http://zeitkapsl.eu/terms" class="text-brand">terms and conditions</a>.',
    "zeitkapsl.accept.privacy.policy":
      'I agree with the <a href="https://zeitkapsl.eu/privacy" class="text-brand">privacy policy</a>.',
    "zeitkapsl.must.accept": "You must agree",

    "form.error.mandatory": " is required",
    "form.error.invalid": " is invalid",

    "invite.code": "Invite code",

    "login.invalid": "Invalid email or password",
    "pwd.invalid": "Password invalid",

    "error.register.password_match": "Your password must match",
    "error.register.invalid_code": "Invalid sign up code",

    "collection.shared.with.me": "Shared with you",
    "shared.collections": "Shared with you",
    "shared.collections.header": "Albums - Shared with you",

    setup: "Set up",
    "authentication.code": "Authentication code",
    activate: "Activate",
    deactivate: "Deactivate",
    verify: "Verify",
    "verify.detail":
      "Open your two-factor authenticator (2FA) app or browser extension to view your authentication code.",
    "settings.2fa.enable.how": "Scan the QR code",
    "settings.2fa.enable.how.detail":
      "Use an authenticator app or browser extension to scan.",
    "settings.2fa.enable.now.you.can.deactivate":
      "Here you can disable two-factor authentication.",
    "settings.2fa.enabled": "Two-factor authentication enabled",
    "settings.2fa.disabled": "Two-factor authentication disabled",
    "settings.access": "Access",
    "settings.access.detail": "Manage your password and authentication (2FA)",
    "settings.2fa": "Two-factor authentication (2FA)",
    "settings.2fa.link": "Enable 2FA <span aria-hidden='true'>&rarr;</span>",
    "settings.2fa.detail":
      "Activate two-factor authentication for extra security using your authenticator app.",
    "settings.2fa.link.deactivate":
      "Disable 2FA <span aria-hidden='true'>&rarr;</span>",
    "settings.language.detail": "Change the language of the user interface.",
    "settings.change_password": "Change password",
    "settings.change_password.detail":
      "Update your password associated with your account.",
    "settings.change_password.current": "Current password",
    "settings.change_password.new": "New password",
    "settings.change_password.link":
      "Change password <span aria-hidden='true'>&rarr;</span>",

    "settings.change_password.new_confirm": "New password confirmation",
    "settings.error.current_password": "Current password is incorrect",
    "settings.password_changed": "Your password was successfully changed",
    "settings.quota.detail":
      "Here you see how much storage space you have available",
    "settings.quota.detail_accounts":
      "Here you can see how many members you can still invite to your subscription.",
    "settings.sessions": "Active devices",
    "settings.sessions.this_device": "this device",
    "settings.sessions.detail": "Here you see your currently active devices",

    "settings.cancel_subscription_and_account.title":
      "Cancel subscription & delete account",
    "settings.cancel_subscription_and_account.detail":
      "You can cancel your subscription and delete your account here",
    "settings.cancel_subscription_and_account.period.end": (periodEnd) =>
      `You can cancel your subscription now. Please note that your subscription will continue until <span class='text-brand font-bold text-lg whitespace-nowrap'>${periodEnd}</span>. Please make sure you have backed up your data by this date!`,
    "settings.cancel_subscription_and_account.description":
      "Please note that your account will be deleted the following day. You will have the opportunity to restore your account and data within 14 days of deletion. Your photos and videos will then be irrevocably deleted.",
    "settings.cancel_subscription_and_account.description2":
      "Your members will be informed of your cancellation. They have the option of taking a subscription themselves. Otherwise, their accounts will be deleted on the same day as yours.",
    "settings.cancel_subscription_and_account": "Cancel subscription",
    "settings.cancel_subscription_and_account.confirmation":
      "Are you sure you want to cancel your subscription and delete your account?",

    "settings.cancel_account": "Delete account",
    "settings.cancel_account.detail":
      "You can request the deletion of your account here",
    "settings.cancel_account.description":
      "Before deleting your account, make sure you have <span class='text-red-600 font-bold text-lg whitespace-nowrap'>backed up</span> your data! Please note that your account will be deleted the following day.",
    "settings.cancel_account.description2":
      "You will have the opportunity to restore your account and data within 14 days of deletion. Your photos and videos will then be irrevocably deleted.",
    "settings.cancel_account.confirmation":
      "Are you sure you want to delete your account?",

    "settings.cancel_account.marked": (cancellation_at) =>
      `Your account will be deleted on <span class='text-red-600 font-bold text-lg whitespace-nowrap'>${cancellation_at}</span>!`,
    "settings.cancel_account.marked.in.days": (days) =>
      `Your account will be deleted in <span class='text-red-600 font-bold text-lg whitespace-nowrap'>${days}</span> day${days > 1 ? "s" : ""}!`,
    "settings.cancel_account.revert": "Ich will mein Konto behalten",
    "settings.cancel_account.revert.confirmation":
      "Bist du dir sicher, dass du dein Konto nicht löschen willst?",

    "settings.cancel_subscription_and_account.buy.abo":
      "Order a subscription <a href='/register/payment' class=\"text-sm font-semibold text-brand cursor-pointer\">here</a> so that your account is not deleted.",

    "settings.subscription": "Current plan",
    "settings.select.subscription": "Select plan",
    "settings.select.subscription.detail": "Here you can upgrade your account.",
    "settings.select.subscription.detail2":
      "You currently don't have a subscription but are a member under someone else. Your account is linked to theirs.",
    "settings.select.subscription.detail3":
      "Subscribe to convert your account into an independent one.",
    "settings.up.to.users": (count) => `Up to ${count} users`,
    "settings.change.subscription":
      "Change plan<span aria-hidden='true'>&rarr;</span>",
    "settings.change.subscription.heading": "Change Plan",
    "settings.select.subscription.heading": "Select Plan",
    "settings.current.subscription": "Current Plan",
    "settings.change.subscription.btn": "Switch",
    "settings.change.subscription.btn.m": "Switch to monthly",
    "settings.change.subscription.btn.y": "Switch to yearly",
    "settings.change.subscription.not.possible": "Switch Not Possible",
    "settings.change.subscription.storageExceeded.detail":
      "You are using more storage than this plan provides.",
    "settings.change.subscription.accountCountExceeded.detail":
      "You are using more users than this plan allows.",
    "pricing.monthly": "Monthly",
    "pricing.yearly": "Annually",
    recommended: "Popular",

    "product.feature.all.included": "All features included",
    "product.feature.always.up.downgrade": "Upgrade anytime",
    "product.feature.cancellable.m": "Cancel anytime, monthly",
    "product.feature.cancellable.y": "Cancel anytime, yearly",

    "reactivate.account": "Reactivate your account",
    "reactivate.account.details":
      "Your account has been deleted. You can recover it by re-activating it here.",
    "reactivate.account.btn": "Reactivate",
    "reactivate.by.checkout.account.details":
      "Your account has been deleted. However, you can recover it by taking a new subscription. </br></br>You will receive the same subscription you had before. If you wish, you can change it at any time in the settings after completion.",
    "reactivate.by.checkout.btn": "Renew subscription",
    "reactivate.by.upgrade.to.paying.account.details":
      "Your account has been deleted. However, you can recover it by taking a subscription.",
    "reactivate.by.upgrade.to.paying.account.btn": "Choose subscription",

    "settings.change_email": "Change email address",
    "settings.change_email.new": "New email address",
    "settings.email_change_requested": "You'll receive an email",
    "settings.email_change_requested.detail": "Confirm by clicking the link",
    "settings.name_changed": "Your name was successfully changed",
    "settings.change_payment_information": "New payment method",

    "settings.invite_to_account": "Invite members",
    "settings.invite_to_account.detail":
      "Share your subscription with your family and friends.",
    "settings.invite_to_account.description": (available) =>
      `There are still <span class='font-bold text-brand'>${available}</span> free ${available > 1 ? "places" : "place"}. Add new members to your account.`,
    "settings.invite_to_account.sent_to": "Send invitation to",
    "settings.invite_to_account.no.accounts.available":
      "You can no longer invite new members. All your accesses are taken.",
    "settings.invite_to_account.send": "Send invite",
    "settings.invite_to_account.sent": "Invite sent",
    "settings.member.of": (name) => `You use ${name}'s subscription.`,

    "settings.members.storage.detail":
      "This is the storage space used by each member.",

    "payment.information": "Payment information",
    "payment.information.details":
      "Here you can manage your payment method and download your recent invoices",
    "payment.information.change": "Manage payment information",
    "payment.information.changed": "Payment information changed",
    "card.ending.with": "Ending with",
    "card.expires": "Expires",

    memories: "Photos",
    favorites: "Favorites",
    trash: "Trash",
    upload: "Upload",
    alert: "Warning!",

    "search.filter": "Filter",
    "search.filter.photos": "Photos",
    "search.filter.videos": "Videos",
    "search.filter.favorites": "Favorites",
    "search.recommendations": "Suggestions",
    "search.recommendations.detail": "2024 or Mai 2024 or 13-01-2024",

    "hero.heading": "simple. secure. private.",
    "hero.heading.part1": "Your photos and videos ",
    "hero.heading.part2": "secure",
    "hero.heading.part3": " from big tech and hackers",
    "hero.subheading": "Privacy by end-to-end encryption",

    download: "Download",
    "download.all": "Download all",
    share: "Share",
    delete: "Delete",
    close: "Close",
    quota: "Storage quota",
    available: "available",
    month: "month",
    year: "year",

    "diff.ago": "ago",
    "diff.seconds": "seconds",
    "diff.hours": "hours",
    "diff.minutes": "minutes",
    "diff.just_now": "just now",
    "diff.days": "days",
    "diff.months": "months",

    purge: "Delete permanently",
    restore: "Restore",

    pricing: "Pricing",

    "landing.features.heading": "Private photos should stay private.",
    "landing.features.heading1": "Private photos",
    "landing.features.heading2": "should stay private.",
    "landing.features.paragraph":
      "We didn't want to leave our photos with the tech giants, so we built an alternative to show there is a different way.",

    "landing.feature.easy.heading": "Easy to use",
    "landing.feature.easy.paragraph":
      "Automatic backups for photos and videos. <br/> In original quality. Across all devices.",

    "landing.feature.secure.heading": "Secure",
    "landing.feature.secure.paragraph":
      "Your photos and videos are end-to-end encrypted using the highest cryptographic standards.",

    "landing.feature.private.heading": "Private",
    "landing.feature.private.paragraph":
      "Only you and those you share the key with can see the photos. Not even us.",

    "landing.feature.trusted.heading": "Trusted",
    "landing.feature.trusted.paragraph":
      "Audited by independent cryptography experts.",

    "landing.feature.social.heading": "Social",
    "landing.feature.social.paragraph":
      "It's easy to share your photos and videos with friends and family. Connected by collaborative albums and links.",

    "landing.feature.eu.heading": "Made in Europe",
    "landing.feature.eu.paragraph":
      "It's developed and operated in Europe and is subject to the GDPR and Austrian privacy laws.",

    collections: "Albums",
    "collections.new": "New album",

    "delete.detail":
      "Remove from your account, devices, and locations where the item was shared?",
    "delete.ok": "Move to trash",

    "purge.detail":
      "Items will not occupy storage space in your account afterwards. Delete permanently?",

    "trash.empty": "No elements",
    "trash.empty.detail":
      "Only photos that have been deleted are displayed here. They will be removed after 60 days.",
    "trash.empty.trash": "Empty trash",
    "trash.empty.trash.header": "Delete permanently?",
    "trash.empty.trash.detail":
      "All items will be permanently deleted. This action can not be undone!",

    "library.empty": "Your library is empty",
    "library.empty.upload": "Do you want to upload anyhting?",
    "library.search.empty": "No results found",

    "upload.drag.and.drop.detail": "Drop the files anywhere to upload them.",

    "pair.new.device.title": "Pair new device",
    "pair.new.device.detail": "Use zeitkapsl on other devices",
    "pair.new.device.btn": "Add device",
    "pair.new.device.description":
      "Open the app on the other device and scan the QR code",

    "favorite.set": "Add to favorites",
    "favorite.unset": "Remove from favorites",

    "months.abbreviations": [
      "jan",
      "feb",
      "mar",
      "apr",
      "may",
      "jun",
      "jul",
      "aug",
      "sep",
      "oct",
      "nov",
      "dec",
    ],

    "months.short": [
      "Jan.",
      "Feb.",
      "Mar.",
      "Apr.",
      "May",
      "Jun.",
      "Jul.",
      "Aug.",
      "Sep.",
      "Oct.",
      "Nov.",
      "Dec.",
    ],

    "months.long": [
      "Jannuary",
      "February",
      "March",
      "April",
      "May",
      "June",
      "July",
      "August",
      "September",
      "October",
      "November",
      "December",
    ],

    "photo.editor.rotate.left": "Rotate left",
    "photo.editor.rotate.right": "Rotate right",
    "photo.editor.flip": "Mirror",
    "photo.editor.revert": "Revert",
    "photo.editor.crop": "Crop",

    "newsletter.receive": "Newsletter",
    subscription_change: "Scheduled subscription change",
    "subscription_change.scheduled_at": "scheduled at",
    "subscription_change.description":
      "Please make sure you reduce your used storage and accounts to match the new package, otherwise the product change will fail. We will remind you again two days before the change will be attempted.",
    "subscription_change.cancel": "Cancel subscription change",
    "subscription_change.confirmation":
      "Do you really want to cancel the subscription change?",
    "subscription_change.confirm": "Change subscription",
    "subscription_change.instant": "Executed immediately",
    "subscription_change.next_cancellation": "Next possible cancellation date",
    "product.select.change": "Change product",
    "retry.all.errors": "Retry all failed",
  },
  "de-DE": {
    "en-US": "EN-US",
    "de-DE": "DE-DE",
    "de-AT": "DE-AT",
    "zeitkapsl.rights": `&copy; ${new Date().getFullYear()} zeitkapsl. Alle Rechte vorbehalten.`,
    login: "Anmelden",
    logout: "Abmelden",
    register: "Registrieren",
    "password.forgot": "Passwort vergessen?",
    email: "E-Mail Adresse",
    password: "Passwort",
    save: "Speichern",
    saving: "Speichere",
    next: "Weiter",
    repeat: "wiederholen",
    "not a member": "Du hast noch keinen Account?",
    settings: "Einstellungen",
    language: "Sprache",
    selected: "ausgewählt",
    unselect: "Selektion aufheben",
    "register.success": "Bitte prüfen deinen E-Mail Posteingang",
    "step.upper.case": "SCHRITT",
    "of.upper.case": "VON",
    of: "von",
    first_name: "Vorname",
    last_name: "Nachname",
    name: "Name",
    "back.to.home": "Zurück",
    back: "Zurück",
    "error.short": "Fehler",
    error: "Ein unerwarteter Fehler ist aufgetreten",
    "try.later": "Bitte versuchen es etwas später!",
    "currently.not.available": "Derzeit nicht verfügbar",
    edit: "Editieren",
    updated: "Aktualisiert",
    conflict: "Das Element wurde verändert",
    "conflict.detail": "Wir haben es neu geladen!",
    elements: "Elemente",
    element: "Element",
    open: "Öffnen",
    create: "Erstellen",
    creating: "Erstelle",
    cancel: "Abbrechen",
    cancelling: "Stoppe",
    tasks: "Aufgaben",
    "tasks.in.progress": "Aufgaben werden noch verarbeitet",
    processing: "Verarbeitung",
    "processing.done": "Verarbeitung abgeschlossen",
    "processing.done.with.errors": "Verarbeitung mit Fehlern abgeschlossen",
    search: "Suche",
    menu: "Menü",
    info: "Info",
    remove: "Entfernen",
    select: "Auswählen",
    "media.not.found": "Media Element nicht gefunden",
    "drag.and.drop.or.click": "Drag und drop oder klicke hier",
    "or.get.the.app": "oder lade die App herunter",
    change: "Ändern",
    "thank.you": "Danke",
    "lets.go": "Los gehts!",
    retry: "Wiederholen",
    account: "Konto",
    accounts: "Konten",
    "email.already.exists": "Email-Adresse bereits vergeben",
    "personal.information": "Persönliche Daten",
    members: "Mitglieder",
    "not.editable": "Schreibgeschützt",
    add: "Hinzufügen",
    "click.here": "Hier klicken",
    "saving.changes": "Änderungen speichern",
    "changes.saved": "Änderungen gespeichert",
    syncing: "Synchronisierung läuft",
    plan: "Plan",
    photos: "Fotos",
    minimize: "Fenster minimieren",
    maximize: "Fenster maximieren",

    "collection.empty": "Dein Album ist leer",
    "collection.empty.add": "Füge etwas deinem Album hinzu",
    "collection.add": "Zu Album hinzufügen",
    "collection.add.to": (collectionName) =>
      `Zum Album '${collectionName}' hinzufügen`,
    "collection.added": "Zu Album hinzugefügt",
    "collection.name": "Albumname",
    "collection.name.enter": "Hier Albumnamen vergeben",
    "collection.delete": "Album löschen",
    "collection.delete.detail":
      "Das Löschen des Albums kann nicht rückgängig gemacht werden. Fotos und Videos aus dem gelöschten Album bleiben in der zeitkapsl gespeichert.",
    "collection.remove": "Aus Album entfernen",
    "collection.remove.detail":
      "Fotos und Videos bleiben weiter in deiner Fotogalerie gespeichert.",
    "collection.goto": "Zum Album",
    "collection.as.cover.media": "Als Albumdeckblatt auswählen",
    "no.collections": "Keine Alben gefunden",
    "no.collections.add": "Möchtest du eines erstellen?",
    "no.shared.collections": "Keine geteilten Alben gefunden",

    "nothing.found": "Keine gefunden",

    "read.only.account": "Demo-Konto",
    "read.only.account.detail": "Daher ist diese Funktion nicht verfügbar.",

    "storage.limit.reached": "Speicherplatz voll",
    "storage.limit.reached.detail":
      "Dein Speicherplatz ist aufgebraucht. Bitte erweitere deinen Speicher.",

    "share.via.link": "Über Link teilen",
    "share.link.create": "Link erstellen",
    "accept.share": "Dieses geteilte Album in meinen Alben merken?",
    "copy.shared.collection": "Album kopieren",
    "copy.shared.collection.success": "Album kopiert",
    "exit.shared.collection": "Entfernen",
    "exit.shared.collection.deleted": "Entfernt",
    "share.link.valid": (validUntil) =>
      `Linkfreigabe, gültig bis ${validUntil}`,
    "share.link.edit": "Linkfreigabe",
    "share.link.allow.uploads": "Uploads zulassen",
    "valid.until": "Gültig bis",
    "share.delete.titel": "Linkfreigabe löschen",
    "share.delete.detail":
      "Das Löschen der Linkfreigabe kann nicht rückgängig gemacht werden. Möchtest Du die Linkfreigabe löschen?",

    waitlist: "In Warteliste eintragen",
    "waitlist.motivation":
      "Sei einer der Ersten, die Zugang zur zeitkapsl erhalten. <br\> <br\> Wenn du niemanden kennst, der eine Einladung hat, kannst du dich auf die Warteliste eintragen und wir senden dir demnächst eine zu.",
    "waitlist.register": "Jetzt anmelden",
    "waitlist.success.title": "Herzlichen Dank!",
    "waitlist.success.detail": "Wir haben dich in der Warteliste hinzugefügt.",

    "register.welcome": "Wähle dein perfektes Abo",
    "register.select.subscription.motivation":
      "Starte jetzt und entdecke, was möglich ist!",
    "register.first.step.motivation":
      "Nur ein paar Schritte und Du hast es geschafft!",
    "register.create.password": "Erstelle dein Konto",
    "register.invite.code": "Einladungscode eingeben",
    "register.invite.code.motivation":
      "Um teilzunehmen benötigst du einen Einladungscode!",
    "register.invite.code.question":
      "Du hast keinen Einladungscode? Trag dich in der ",
    "register.invite.code.question.resolve": "Warteliste",
    "register.invite.code.question.resolve.final": " ein.",
    "register.user.password.strength": "Passwortstärke: ",
    "register.user.password.strength.too.weak": "sehr schwach",
    "register.user.password.strength.weak": "schwach",
    "register.user.password.strength.medium": "mittel",
    "register.user.password.strength.strong": "stark",
    "register.payment.information": "Wähle deine bevorzugte Zahlungsart",
    "register.payment.information.details1":
      "Deine Zahlungen sind bei uns sicher. Du kannst bequem online kündigen.",
    "register.payment.information.details2":
      "Deine Angaben sind verschlüsselt und du kannst die Zahlungsart jederzeit ändern.",
    "register.recovery": "Speichere deinen Schlüssel",
    "register.recovery.details":
      "Dein Schlüssel ist geheim. Mit diesem Schlüssel werden deine Daten verschlüsselt. Wir speichern deinen Schlüssel <span class='font-bold'>nicht!</span> <br\> <br\> Wenn du dein Passwort vergisst, kannst du deine Daten nur mit diesem Schlüssel wiederherstellen. <br\> <br\> Lade den Schlüssel herunter und bewahre ihn an einem sicheren Ort auf.",
    "register.recovery.details.short":
      "Dein Schlüssel ist geheim. Mit diesem Schlüssel werden deine Daten verschlüsselt. Wir speichern deinen Schlüssel <span class='font-bold'>nicht!</span> <br\> <br\> Wenn du dein Passwort vergisst, kannst du deine Daten nur mit diesem Schlüssel wiederherstellen.",
    "register.recovery.wordlist":
      "Deine persönlichen Wiederherstellungs-Wörter",
    "register.recovery.qrcode": "QR-Code",
    "register.recovery.qrcode.detail":
      "Zum Scannen beim Zurücksetzen des Passwords.",
    "register.recovery.secret.key": "Schlüssel",

    "register.recovery.pdf.title": "Emergency Kit",
    "register.recovery.pdf.description":
      "Wenn du dein Passwort vergisst, kannst du deine Daten nur mit diesem Schlüssel wiederherstellen und dein Passwort zurücksetzen.** Wir können weder auf dein Daten zugreifen noch diese wiederherstellen.**\n\nBewahre deinen Schlüssel an einem sicheren Ort auf, beispielsweise bei deiner Geburtsurkunde.",
    "register.recovery.pdf.account": "Erstellt für",
    "register.recovery.needhelp": "Brauchst du Hilfe?",
    "register.recovery.needhelp.detail": "support@zeitkapsl.eu",
    "register.recovery.pdf.download": "Geheimen Schlüssel herunterladen",
    "register.recovery.pdf.download.again": "Nochmals herunterladen",
    "register.in_progress": "Wir legen deinen Account an",
    "register.for_free": "Für dich ist das gratis :)",
    "register.now": "Jetzt registrieren",
    "register.payment.success.detail":
      "Vielen Dank für Deine Zahlung! Dein Konto ist bereit.",

    "payment.failed": "Zahlung fehlgeschlagen",

    "recovery.key": "Schlüssel",
    "recovery.scan.qr.code": "Scan QR-Code",
    "recovery.scan.qr.code.stop": "Scanner schließen",

    "reset.password": "Passwort zurücksetzen",
    "reset.password.detail":
      "Gib deine E-Mail Adresse und ein neues Passwort ein, um das Passwort für dein Konto zurückzusetzen.",
    "reset.password.new": "Neues Passwort eingeben",
    "reset.password.do": "Passwort zurücksetzen",

    "zeitkapsl.accept.terms.and.conditions":
      'Ich habe die <a href="http://zeitkapsl.eu/terms" class="text-brand">allgemeinen Geschäftsbedingungen</a> gelesen und stimme diesen zu.',
    "zeitkapsl.accept.privacy.policy":
      'Ich akzeptiere die <a href="https://zeitkapsl.eu/privacy" class="text-brand">Datenschutzbestimmungen</a>.',
    "zeitkapsl.must.accept": "Zustimmung erforderlich",

    "form.error.mandatory": " ist erforderlich",
    "form.error.invalid": " ist ungültig",

    "login.invalid": "Email oder Passwort ungültig",
    "pwd.invalid": "Passwort ungültig",

    "invite.code": "Einladungscode",

    "error.register.password_match": "Dein Passwort muss übereinstimmen",
    "error.register.invalid_code": "Ungültiger Registrierungscode",

    "collection.shared.with.me": "Mit dir geteilt",
    "shared.collections": "Geteilt mit dir",
    "shared.collections.header": "Alben - Geteilt mit dir",

    setup: "Einrichten",
    "authentication.code": "Authentifizierungs-Code",
    activate: "Aktivieren",
    deactivate: "Deaktivieren",
    verify: "Verifizieren",
    "verify.detail":
      "Öffne deine Zwei-Faktor-Authentifizierungs-App (2FA) oder Browser-Erweiterung, um deinen Authentifizierungscode anzuzeigen.",
    "settings.2fa.enable.how": "Scanne den QR-Code",
    "settings.2fa.enable.how.detail":
      "Verwende eine Authentifizierungs-App oder eine Browser-Erweiterung zum Scannen.",
    "settings.2fa.enable.now.you.can.deactivate":
      "Hier kannst du die Zwei-Faktor-Authentifizierung deaktivieren.",
    "settings.2fa.enabled": "Zwei-Faktor-Authentifizierung aktiviert",
    "settings.2fa.disabled": "Zwei-Faktor-Authentifizierung deaktiviert",
    "settings.access": "Zugriff",
    "settings.access.detail":
      "Verwalte dein Passwort und deine Authentifizierung (2FA)",
    "settings.2fa": "Zwei-Faktor-Authentifizierung (2FA)",
    "settings.2fa.link": "2FA aktivieren<span aria-hidden='true'>&rarr;</span>",
    "settings.2fa.detail":
      "Aktiviere die Zwei-Faktor-Authentifizierung für zusätzliche Sicherheit mit deiner Authentifizierungs-App.",
    "settings.2fa.link.deactivate":
      "2FA deaktivieren<span aria-hidden='true'>&rarr;</span>",
    "settings.language.detail": "Ändere die Sprache der Benutzeroberfläche.",
    "settings.change_password": "Passwort ändern",
    "settings.change_password.link":
      "Passwort ändern <span aria-hidden='true'>&rarr;</span>",
    "settings.change_password.detail":
      "Ändere das Passwort das mit deinem Account verknüpft ist.",
    "settings.change_password.current": "Aktuelles Passwort",
    "settings.change_password.new": "Neues Passwort",
    "settings.change_password.new_confirm": "Neues Passwort wiederholen",
    "settings.error.current_password": "Das aktuelle Passwort ist inkorrekt",
    "settings.password_changed": "Dein Passwort wurde erfolgreich geändert",
    "settings.quota.detail":
      "Hier siehst du wie viel Speicherplatz du zur Verfügung hast",
    "settings.quota.detail_accounts":
      "Hier siehst du wie viele Mitglieder du zu deinem Abo noch einladen kannst.",
    "settings.sessions": "Aktive Geräte",
    "settings.sessions.this_device": "dieses Gerät",
    "settings.sessions.detail":
      "Hier siehst du deine angemeldeten aktiven Geräte",

    "settings.cancel_subscription_and_account.title":
      "Abo kündigen & Konto löschen",
    "settings.cancel_subscription_and_account.detail":
      "Hier kannst du dein Abo kündigen und dein Konto löschen",
    "settings.cancel_subscription_and_account.period.end": (periodEnd) =>
      `Du kannst dein Abo jetzt kündigen. Beachte, dass dein Abo dann noch bis zum <span class='text-brand font-bold text-lg whitespace-nowrap'>${periodEnd}</span> weiterläuft. Bitte stelle sicher, dass du deine Daten bis zu diesem Datum gesichert hast!`,
    "settings.cancel_subscription_and_account.description":
      "Beachte, dass dein Konto am darauf folgenden Tag gelöscht wird. Innerhalb der nächsten 14 Tage nach der Löschung besteht die Möglichkeit, dein Konto und deine Daten wiederherzustellen. Danach werden deine Fotos und Videos unwiderruflich gelöscht.",
    "settings.cancel_subscription_and_account.description2":
      "Deine Mitglieder werden über deine Kündigung informiert. Sie haben die Möglichkeit selbst ein Abonnement abzuschließen. Andernfalls werden ihre Konten am gleichen Tag gelöscht wie deins.",
    "settings.cancel_subscription_and_account": "Abo kündigen",
    "settings.cancel_subscription_and_account.confirmation":
      "Bist du dir sicher, dass du dein Abo kündigen und dein Konto löschen willst?",

    "settings.cancel_account": "Konto löschen",
    "settings.cancel_account.detail":
      "Hier kannst du die Löschung deines Kontos veranlassen",
    "settings.cancel_account.description":
      "Bevor du dein Konto löschst, vergewissere dich, dass du deine Daten <span class='text-red-600 font-bold text-lg whitespace-nowrap'>gesichert</span> hast! Beachte, dass dein Konto am folgenden Tag gelöscht wird.",
    "settings.cancel_account.description2":
      "Innerhalb der nächsten 14 Tage nach der Löschung besteht die Möglichkeit, dein Konto und deine Daten wiederherzustellen. Danach werden deine Fotos und Videos unwiderruflich gelöscht.",
    "settings.cancel_account.confirmation":
      "Bist du dir sicher, dass du dein Konto löschen willst?",

    "settings.cancel_account.marked": (cancellation_at) =>
      `Dein Konto wird am <span class='text-red-600 font-bold text-lg whitespace-nowrap'>${cancellation_at}</span> gelöscht!`,
    "settings.cancel_account.marked.in.days": (days) =>
      `Dein Konto wird in <span class='text-red-600 font-bold text-lg whitespace-nowrap'>${days}</span> Tag${days > 1 ? "en" : ""} gelöscht!`,
    "settings.cancel_account.revert": "Konto nicht löschen",
    "settings.cancel_account.revert.confirmation":
      "Bist du dir sicher, dass du dein Konto nicht löschen willst?",

    "settings.cancel_subscription_and_account.buy.abo":
      "Bestelle <a href='/register/payment' class=\"text-sm font-semibold text-brand cursor-pointer\">hier</a> ein Abo, damit dein Konto nicht gelöscht wird.",

    "settings.subscription": "Aktuelles Abo",
    "settings.select.subscription": "Abo auswählen",
    "settings.select.subscription.detail": "Du kannst dein Konto upgraden.",
    "settings.select.subscription.detail2":
      "Du besitzt derzeit kein Abo, sondern bist bei jemanden Mitglied. Dein Konto ist mit diesem verknüpft.",
    "settings.select.subscription.detail3":
      "Schließe ein Abo ab, um dein Konto in ein eigenständiges umzuwandeln.",
    "settings.up.to.users": (count) => `Bis zu ${count} User`,
    "settings.change.subscription":
      "Abo ändern<span aria-hidden='true'>&rarr;</span>",
    "settings.change.subscription.heading": "Abo ändern",
    "settings.select.subscription.heading": "Abo wählen",
    "settings.current.subscription": "Aktuelles Abo",
    "settings.change.subscription.btn": "Wechseln",
    "settings.change.subscription.btn.m": 'Zu "monatlich" wechseln',
    "settings.change.subscription.btn.y": 'Zu "jährlich" wechseln',
    "settings.change.subscription.not.possible": "Wechsel nicht möglich",
    "settings.change.subscription.storageExceeded.detail":
      "Du belegst mehr Speicher als dieses Paket bietet.",
    "settings.change.subscription.accountCountExceeded.detail":
      "Du hast mehr User in gebrauch als dieses Paket bietet.",
    "pricing.monthly": "Monatlich",
    "pricing.yearly": "Jährlich",
    recommended: "Beliebt",

    "product.feature.all.included": "Alle Features inkludiert",
    "product.feature.always.up.downgrade": "Upgrade jederzeit möglich",
    "product.feature.cancellable.m": "Monatlich kündbar",
    "product.feature.cancellable.y": "Jährlich kündbar",

    "settings.change_email": "E-Mail Adresse ändern",
    "settings.change_email.detail":
      "Ändere deine E-Mail Adresse die mit deinem Account verknüpft ist.",
    "settings.change_email.new": "Neue E-Mail Adresse",
    "settings.email_change_requested": "Du bekommst eine E-Mail",
    "settings.email_change_requested.detail": "Rufe den Bestätigungslink auf",
    "settings.name_changed": "Dein Name wurde erfolgreich geändert",
    "settings.change_payment_information": "Neue Zahlungsmethode erfassen",

    "settings.invite_to_account": "Mitglieder",
    "settings.invite_to_account.detail":
      "Nutze dein Abo gemeinsam mit deiner Familie und Freunden.",
    "settings.invite_to_account.description": (available) =>
      `Es gibt noch <span class='font-bold text-brand text-xl'>${available}</span> ${available > 1 ? "freie Plätze" : "freien Platz"}. Füge neue Mitglieder deinem Konto hinzu. Sende eine Einladung.`,
    "settings.invite_to_account.sent_to": "Einladung senden an",
    "settings.invite_to_account.no.accounts.available":
      "Du kannst keine neuen Mitglieder mehr einladen. Alle deine Zugänge sind vergeben.",
    "settings.invite_to_account.send": "Einladung senden",
    "settings.invite_to_account.sent": "Einladung gesendet",
    "settings.member.of": (name) => `Du nutzt das Abo von ${name} mit.`,

    "settings.members.storage.detail":
      "Das ist der durch die einzelnen Mitglieder belegte Speicherplatz.",

    "reactivate.account": "Reaktiviere dein Konto",
    "reactivate.account.details":
      "Dein Konto wurde gelöscht. Du kannst es jedoch wiederherstellen, indem Du es hier reaktivierst.",
    "reactivate.account.btn": "Reaktivieren",
    "reactivate.by.checkout.account.details":
      "Dein Konto wurde gelöscht. Du kannst es jedoch wiederherstellen, indem du ein neues Abonnement abschließt. </br></br>Dabei erhältst du das gleiche Abo, das du zuvor hattest. Falls du möchtest, kannst du es nach Abschluss jederzeit in den Einstellungen ändern.",
    "reactivate.by.checkout.btn": "Abonnement erneuern",
    "reactivate.by.upgrade.to.paying.account.details":
      "Dein Konto wurde gelöscht. Du kannst es jedoch wiederherstellen, indem du ein neues Abonnement abschließt.",
    "reactivate.by.upgrade.to.paying.account.btn": "Abonnement auswählen",

    "payment.information": "Zahlungsinformationen",
    "payment.information.details":
      "Hier kannst du deine Zahlungsmethode verwalten und deine letzten Rechnungen herunterladen",
    "payment.information.change": "Zahlungsinformationen verwalten",
    "payment.information.changed": "Zahlungsinformationen geändert",
    "card.ending.with": "Endet mit",
    "card.expires": "Läuft ab",

    memories: "Fotos",
    favorites: "Favoriten",
    trash: "Papierkorb",
    upload: "Hochladen",
    alert: "Achtung!",

    "search.filter": "Filter",
    "search.filter.photos": "Fotos",
    "search.filter.videos": "Videos",
    "search.filter.favorites": "Favoriten",
    "search.recommendations": "Vorschläge",
    "search.recommendations.detail": "2024 oder Mai 2024 oder 13.05.2024",

    "hero.heading.part1": "Deine Fotos und Videos ",
    "hero.heading.part2": "sicher",
    "hero.heading.part3": " vor Tech-Konzernen und Hackern",
    "hero.subheading": "Privatsphäre durch Ende-zu-Ende Verschlüsselung",

    download: "Herunterladen",
    "download.all": "Alle herunterladen",
    share: "Teilen",
    delete: "Löschen",
    close: "Schließen",
    quota: "Speicherplatz",
    available: "verfügbar",
    month: "Monat",
    year: "Jahr",

    "diff.ago": "vor",
    "diff.seconds": "Sekunden",
    "diff.hours": "Stunden",
    "diff.minutes": "Minuten",
    "diff.days": "Tage",
    "diff.months": "Monate",
    "diff.just_now": "gerade eben",

    purge: "Endgültig löschen",
    restore: "Wiederherstellen",

    pricing: "Preise",

    "landing.features.heading":
      "Private Fotos und Videos sollten privat bleiben.",
    "landing.features.heading1": "Private photos",
    "landing.features.heading2": "should stay private.",
    "landing.features.paragraph":
      "Wir wollten unsere Fotos nicht bei den großen Tech-Konzernen lassen, also haben wir eine Alternative entwickelt, um zu zeigen, dass es einen anderen Weg gibt.",

    "landing.feature.easy.heading": "Einfach",
    "landing.feature.easy.paragraph":
      "Automatische Backups für Fotos und Videos. In Originalqualität. Auf allen Geräten.",

    "landing.feature.secure.heading": "Sicher",
    "landing.feature.secure.paragraph":
      "Deine Fotos und Videos werden Ende-zu-Ende verschlüsselt unter Verwendung der höchsten kryptografischen Standards.",

    "landing.feature.private.heading": "Privat",
    "landing.feature.private.paragraph":
      "Nur du und diejenigen, mit denen du den Schlüssel teilst, können die Fotos sehen. Nicht einmal wir.",

    "landing.feature.trusted.heading": "Vertrauenswürdig",
    "landing.feature.trusted.paragraph":
      "Von unabhängigen Sicherheitsspezialisten geprüft.",

    "landing.feature.social.heading": "Teilen",
    "landing.feature.social.paragraph":
      "Fotos und Videos können einfach durch gemeinsame Alben und Links mit Freunden und Familie geteilt werden.",

    "landing.feature.eu.heading": "Made in Europe",
    "landing.feature.eu.paragraph":
      "Die Platfform wird in Europa entwickelt, betrieben und unterliegt der DSGVO sowie den österreichischen Datenschutzgesetzen.",

    collections: "Alben",
    "collections.new": "Neues Album",

    "delete.detail":
      "Aus deinem Konto, von Geräten und an Orten, an denen das Element geteilt wurden, entfernen?",
    "delete.ok": "In den Papierkorb verschieben",

    "purge.detail":
      "Elemente werden danach nicht auf den Speicherplatz in deinem Konto angerechnet. Element endgültig löschen?",

    "trash.empty": "Keine Fotos",
    "trash.empty.detail":
      "Hier werden nur Fotos angezeigt, die gelöscht wurden. Sie werden nach 60 Tagen entfernt.",
    "trash.empty.trash": "Papierkorb leeren",
    "trash.empty.trash.header": "Endgültig löschen?",
    "trash.empty.trash.detail":
      "Alle Elemente werden endgültig gelöscht. Dieser Vorgang kann nicht rückgängig gemacht werden!",

    "library.empty": "Deine Bibliothek ist leer",
    "library.empty.upload": "Möchtest du etwas hochladen?",
    "library.search.empty": "Keine Ergebnisse gefunden",

    "upload.drag.and.drop.detail":
      "Leg die Dateien an einer beliebigen Stelle ab, um sie hochzuladen.",

    "pair.new.device.title": "Gerät hinzufügen",
    "pair.new.device.detail": "Verwende zeitkapsl auf anderen Geräten",
    "pair.new.device.btn": "Gerät hinzufügen",
    "pair.new.device.description":
      "Öffne die App auf dem anderen Gerät und scanne den QR-Code",

    "favorite.set": "Zu Favoriten hinzufügen",
    "favorite.unset": "Aus Favoriten entfernen",

    "months.abbreviations": [
      "jan",
      "feb",
      "mär",
      "apr",
      "mai",
      "jun",
      "jul",
      "aug",
      "sep",
      "okt",
      "nov",
      "dez",
    ],

    "months.short": [
      "Jan.",
      "Feb.",
      "Mär.",
      "Apr.",
      "Mai",
      "Jun.",
      "Jul.",
      "Aug.",
      "Sep.",
      "Okt.",
      "Nov.",
      "Dez.",
    ],

    "months.long": [
      "Januar",
      "Februar",
      "März",
      "April",
      "Mai",
      "Juni",
      "Juli",
      "August",
      "September",
      "Oktober",
      "November",
      "Dezember",
    ],

    "photo.editor.rotate.left": "Links rotieren",
    "photo.editor.rotate.right": "Rechts rotieren",
    "photo.editor.flip": "Spiegeln",
    "photo.editor.revert": "Zurücksetzen",
    "photo.editor.crop": "Zuschneiden",

    "newsletter.receive": "Newsletter",
    subscription_change: "Geplante Abo-Änderung",
    "subscription_change.scheduled_at": "wird durchgeführt am",
    "subscription_change.description":
      "Bitte stelle sicher dass du deinen verwendeten Speicherplatz und die Accounts rechtzeitig reduzierst, andernfalls wird der Paketwechsel nicht durchgeführt. Wir werden dich zwei Tage vor dem Wechsel daran erinnern.",
    "subscription_change.cancel": "Abo-Änderung zurückziehen",
    "subscription_change.confirmation":
      "Möchtest du wirklich den Paketwechsel zurückziehen?",
    "subscription_change.confirm": "Abo ändern",
    "subscription_change.instant": "Wird sofort durchgeführt",
    "subscription_change.next_cancellation": "Nächstmögliche Kündigung am",
    "product.select.change": "Product ändern",
    "retry.all.errors": "Alle fehlgeschlagenen wiederholen",
  },
  "de-AT": {
    "months.abbreviations": [
      "jän",
      "feb",
      "mär",
      "apr",
      "mai",
      "jun",
      "jul",
      "aug",
      "sep",
      "okt",
      "nov",
      "dez",
    ],
    "months.short": [
      "Jän.",
      "Feb.",
      "Mär.",
      "Apr.",
      "Mai",
      "Jun.",
      "Jul.",
      "Aug.",
      "Sep.",
      "Okt.",
      "Nov.",
      "Dez.",
    ],
    "months.long": [
      "Jänner",
      "Februar",
      "März",
      "April",
      "Mai",
      "Juni",
      "Juli",
      "August",
      "September",
      "Oktober",
      "November",
      "Dezember",
    ],
  },
};
