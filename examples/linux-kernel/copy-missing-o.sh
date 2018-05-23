#!/usr/bin/env bash

### Add the object files generated from assembly to the build folder

export build_home=$HOME/standalone-build
export ker=$HOME/linux-stable

#asm files
export lib_obj_path=./built-ins/objects/lib_assembly_objects
export arch_obj_path=./built-ins/objects/arch_assembly_objects
export xlib_obj_path=./built-ins/objects/xlib_assembly_objects
export pow_obj_path=./built-ins/objects/pow_assembly_objects

cp $ker/arch/x86/entry/entry_64.o $arch_obj_path
cp $ker/arch/x86/entry/thunk_64.o $arch_obj_path
cp $ker/arch/x86/entry/vsyscall/vsyscall_emu_64.o $arch_obj_path
cp $ker/arch/x86/entry/entry_64_compat.o $arch_obj_path
cp $ker/arch/x86/realmode/rmpiggy.o $arch_obj_path
cp $ker/arch/x86/kernel/acpi/wakeup_64.o $arch_obj_path
cp $ker/arch/x86/kernel/relocate_kernel_64.o $arch_obj_path
cp $ker/arch/x86/platform/efi/efi_stub_64.o $arch_obj_path

cp $ker/lib/lib-ksyms.o $lib_obj_path

cp $ker/arch/x86/lib/iomap_copy_64.o $xlib_obj_path
cp $ker/arch/x86/lib/hweight.o $xlib_obj_path
cp $ker/arch/x86/lib/msr-reg.o $xlib_obj_path
cp $ker/arch/x86/lib/lib-ksyms.o $xlib_obj_path

cp $ker/arch/x86/power/hibernate_asm_64.o $pow_obj_path

#kernel obj files
export ker_objs_path=./built-ins/objects/ker_objects
cp $ker/arch/x86/kernel/head64.o $ker_objs_path
cp $ker/arch/x86/kernel/head_64.o  $ker_objs_path
cp $ker/arch/x86/kernel/platform-quirks.o $ker_objs_path
cp $ker/arch/x86/kernel/ebda.o $ker_objs_path
cp $ker/usr/initramfs_data.o $ker_objs_path

#arch lib asm
export libx_obj_path=./built-ins/objects/libx_objects
cp $ker/arch/x86/lib/memset_64.o $libx_obj_path
cp $ker/arch/x86/lib/getuser.o $libx_obj_path
cp $ker/arch/x86/lib/rwsem.o $libx_obj_path
cp $ker/arch/x86/lib/memcpy_64.o $libx_obj_path
cp $ker/arch/x86/lib/memmove_64.o $libx_obj_path
cp $ker/arch/x86/lib/copy_user_64.o $libx_obj_path
cp $ker/arch/x86/lib/putuser.o $libx_obj_path
cp $ker/arch/x86/lib/csum-copy_64.o $libx_obj_path
cp $ker/arch/x86/lib/clear_page_64.o $libx_obj_path
cp $ker/arch/x86/lib/copy_page_64.o $libx_obj_path
cp $ker/arch/x86/lib/cmpxchg16b_emu.o $libx_obj_path
cp $ker/arch/x86/lib/retpoline.o $libx_obj_path