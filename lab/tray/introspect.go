package main

const IntrospectXML = `<!DOCTYPE node PUBLIC "-//freedesktop//DTD D-BUS Object Introspection 1.0//EN"
        "http://www.freedesktop.org/standards/dbus/1.0/introspect.dtd">
<node>
    <interface name="org.kde.StatusNotifierWatcher">
        <property name="IsStatusNotifierHostRegistered" type="b" access="read"/>
        <property name="ProtocolVersion" type="i" access="read"/>
        <property name="RegisteredStatusNotifierItems" type="as" access="read"/>
        <signal name="StatusNotifierItemRegistered">
            <arg name="service" type="s" direction="out"/>
        </signal>
        <signal name="StatusNotifierItemUnregistered">
            <arg name="service" type="s" direction="out"/>
        </signal>
        <signal name="StatusNotifierHostRegistered">
        </signal>
        <method name="RegisterStatusNotifierItem">
            <arg name="serviceOrPath" type="s" direction="in"/>
        </method>
        <method name="RegisterStatusNotifierHost">
            <arg name="service" type="s" direction="in"/>
        </method>
    </interface>
    <interface name="org.freedesktop.DBus.Properties">
        <method name="Get">
            <arg name="interface_name" type="s" direction="in"/>
            <arg name="property_name" type="s" direction="in"/>
            <arg name="value" type="v" direction="out"/>
        </method>
        <method name="Set">
            <arg name="interface_name" type="s" direction="in"/>
            <arg name="property_name" type="s" direction="in"/>
            <arg name="value" type="v" direction="in"/>
        </method>
        <method name="GetAll">
            <arg name="interface_name" type="s" direction="in"/>
            <arg name="values" type="a{sv}" direction="out"/>
            <annotation name="org.qtproject.QtDBus.QtTypeName.Out0" value="QVariantMap"/>
        </method>
        <signal name="PropertiesChanged">
            <arg name="interface_name" type="s" direction="out"/>
            <arg name="changed_properties" type="a{sv}" direction="out"/>
            <arg name="invalidated_properties" type="as" direction="out"/>
        </signal>
    </interface>
    <interface name="org.freedesktop.DBus.Introspectable">
        <method name="Introspect">
            <arg name="xml_data" type="s" direction="out"/>
        </method>
    </interface>
    <interface name="org.freedesktop.DBus.Peer">
        <method name="Ping"/>
        <method name="GetMachineId">
            <arg name="machine_uuid" type="s" direction="out"/>
        </method>
    </interface>
</node>`
