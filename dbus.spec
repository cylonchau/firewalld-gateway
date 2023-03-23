%define name dbus
%define version 1.12.20
%define release 1
%define release_tag %{?dist}

Name:           %{name}
Version:        %{version}
Release:        %{release}%{release_tag}
Summary:        D-Bus message bus system

Group:          System Environment/Daemons
License:        AFL-2.1 or GPL-2.0 or GPL-3.0
URL:            https://dbus.freedesktop.org/
Source0:        https://dbus.freedesktop.org/releases/%{name}/%{name}-%{version}.tar.gz

BuildRequires:  gcc
BuildRequires:  make
BuildRequires:  pkgconfig
BuildRequires:  xmlto

%description
D-Bus is a message bus system, a simple way for applications to talk to one another. In addition to interprocess communication, D-Bus helps coordinate process lifecycle; it makes it simple and reliable to code a "single instance" application or daemon, and to launch applications and daemons on demand when their services are needed.

%package libs
Summary:        Libraries for D-Bus
Group:          System Environment/Libraries
License:        AFL-2.1 or GPL-2.0 or GPL-3.0
Requires:       %{name} = %{version}-%{release}%{?dist}

%description libs
This package contains libraries needed at runtime for applications using D-Bus.
%prep
%autosetup -p1

%build
%configure
make %{?_smp_mflags}

%install
%make_install


%check
make check

%clean
rm -rf %{buildroot}

%post -p /sbin/ldconfig
%postun -p /sbin/ldconfig

%files libs
%defattr(-,root,root,-)
%{_libdir}/libdbus-1.*


%files
%{_sysconfdir}/*
%{_bindir}/dbus-cleanup-sockets
%{_bindir}/dbus-daemon
%{_bindir}/dbus-launch
%{_bindir}/dbus-monitor
%{_bindir}/dbus-run-session
%{_bindir}/dbus-send
%{_bindir}/dbus-uuidgen
%{_mandir}/man1/*
%{_datadir}/dbus-1/
%{_prefix}/lib/systemd/system/*
%{_prefix}/libexec/*

%changelog
* Fri Mar 21 2023 John Doe <cylonchau@outlook.com> - 1.12.20-1
- Initial build of D-Bus 1.12.20
